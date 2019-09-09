package guardmech

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var randLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var clientID string
var clientSecret string
var redirectURL string

const authSessionKey string = `_guardmech_csrf`
const sessionKey string = `_guardmech_session`

func init() {
	clientID = os.Getenv("GUARDMECH_CLIENT_ID")
	clientSecret = os.Getenv("GUARDMECH_CLIENT_SECRET")
	redirectURL = os.Getenv("GUARDMECH_REDIRECT_URL")
}

type AuthMux struct {
	conf     *oauth2.Config
	provider *oidc.Provider
}

func generateRandomString(length int, letters []rune) string {
	b := strings.Builder{}
	b.Grow(length)
	lc := len(letters)

	for i := 0; i < length; i++ {
		b.WriteRune(letters[rand.Intn(lc)])
	}
	return b.String()
}

func NewAuthMux() *AuthMux {
	ctx := context.Background()
	p, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		// handle error
		panic(`failed to initialize Google OpenID Connect provider:` + err.Error())
	}

	am := &AuthMux{
		conf: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     p.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		},
		provider: p,
	}

	return am
}

func (a *AuthMux) Mux() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/auth/start", a.StartAuth)
	m.HandleFunc("/auth/callback", a.CallbackAuth)
	m.HandleFunc("/auth/request", a.AuthRequest)
	m.HandleFunc("/auth/sign_in", a.NeedAuth)

	return m
}

func (a *AuthMux) StartAuth(w http.ResponseWriter, req *http.Request) {
	// return path
	if req.ParseForm() != nil {
		http.Error(w, "Body Parsing Error", http.StatusBadRequest)
		return
	}
	path := req.Form.Get("p")
	if path == "" {
		path = "%2F"
	}

	csrfToken := generateRandomString(32, randLetters)
	as := &AuthSession{
		CSRFToken: csrfToken,
		Path:      path,
	}

	authSession := &SessionPayload{
		ExpireAt: time.Now().Add(5 * time.Minute),
		Data:     as,
	}

	// bake CSRF cookie
	http.SetCookie(w, authSession.MakeCookie(req, authSessionKey, 0))
	http.Redirect(w, req, a.conf.AuthCodeURL(csrfToken), http.StatusFound)
}

func (a *AuthMux) CallbackAuth(w http.ResponseWriter, req *http.Request) {
	// session validation
	c, err := req.Cookie(authSessionKey)
	if err != nil {
		http.Error(w, "No CSRF Session", http.StatusForbidden)
		return
	}

	as := &AuthSession{}
	authSession, err := NewSessionPayload(c.Value, as)
	if err != nil {
		log.Println(err)
		http.Error(w, "Session Validation Failed", http.StatusForbidden)
	}
	// delete cookie
	http.SetCookie(w, authSession.RevokeCookie(req, authSessionKey))

	// decode path
	path, err := url.PathUnescape(as.Path)
	if err != nil {
		path = "/"
	}

	// CSRF check
	req.ParseForm()
	state := req.Form.Get("state")
	if state != as.CSRFToken {
		http.Error(w, "Detect Possibility of CSRF Attack", http.StatusForbidden)
		return
	}

	// ID Token
	code := req.URL.Query().Get("code")
	token, err := a.verifyAuthCode(req.Context(), code)
	if err != nil {
		http.Error(w, "Failed to Verify AuthCode", http.StatusForbidden)
		return
	}

	// normalize issuer (remove "https://" for Google)
	token.Issuer = strings.Replace(token.Issuer, "https://", "", -1)

	// if first user, set as owner
	ctx := req.Context()
	conn, err := db.Conn(ctx)

	ok, err := HasPrincipal(ctx, conn)
	if err != nil {
		http.Error(w, "Failed to Check Principal", http.StatusInternalServerError)
		return
	}
	// first cut
	if !ok {
		log.Println("No User!! Entering setup mode.")

		err := CreateFirstPrincipal(ctx, conn, token)
		if err != nil {
			log.Println("failed to setup:", err.Error())
			http.Error(w, "Failed to Setup", http.StatusInternalServerError)
			return
		}
	}

	is := &IDSession{
		Issuer:  token.Issuer,
		Subject: token.Sub,
		Email:   token.Email,
	}

	// OK!! Create Session!!
	session := &SessionPayload{
		Data:     is,
		ExpireAt: time.Now().Add(5 * time.Minute),
	}
	http.SetCookie(w, session.MakeCookie(req, sessionKey, 24*time.Hour))

	http.Redirect(w, req, path, http.StatusFound)
}

type OpenIDToken struct {
	Issuer   string `json:"iss"`
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
	Name     string `json:"name"`
}

func (a *AuthMux) verifyAuthCode(ctx context.Context, code string) (*OpenIDToken, error) {
	var verifier = a.provider.Verifier(&oidc.Config{ClientID: a.conf.ClientID})
	oauth2Token, err := a.conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("Token does not contains id_token")
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.Wrap(err, "Verification failed")
	}

	// Extract custom claims
	var claims OpenIDToken
	if err := idToken.Claims(&claims); err != nil {
		// handle error
		return nil, err
	}

	return &claims, nil
}

func (a *AuthMux) parseIDSession(req *http.Request, is *IDSession) (*SessionPayload, error) {
	// parse session cookie
	c, err := req.Cookie(sessionKey)
	if err != nil {
		return nil, fmt.Errorf("first visit")
	}

	session, err := NewSessionPayload(c.Value, is)
	if err != nil {
		return nil, fmt.Errorf("invalid cookie")
	}

	return session, nil
}

func (a *AuthMux) AuthRequest(w http.ResponseWriter, req *http.Request) {
	is := &IDSession{}
	session, err := a.parseIDSession(req, is)
	if err != nil {
		//log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check ExpireAt
	if time.Now().Sub(session.ExpireAt) > 0 {
		//log.Println("expired")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// fetch Role & Group
	ctx := req.Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Println("Could Not Get Conn")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pri, err := FindPrincipal(ctx, conn, is.Issuer, is.Subject)
	if err != nil {
		log.Println("Could Not Find Principal: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload, err := FindPrincipalByID(ctx, conn, fmt.Sprintf("%d", pri.ID))
	if err != nil {
		log.Println("Could Not Find PrincipalPayload: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// OK! print headers
	w.Header().Set("X-Guardmech-Email", is.Email)
	if len(payload.Groups) > 0 {
		w.Header().Set("X-Guardmech-Groups", payload.Groups[0].Name)
	}
	if len(payload.Roles) > 0 {
		w.Header().Set("X-Guardmech-Roles", payload.Roles[0].Name)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *AuthMux) NeedAuth(w http.ResponseWriter, req *http.Request) {
	// catch path
	originUri := req.Header.Get("X-Auth-Request-Redirect")
	path := url.PathEscape(originUri)

	// cookie check
	is := &IDSession{}
	session, err := a.parseIDSession(req, is)
	if err == nil {
		// cookie is live, try to authenticate without prompt
		if time.Now().Sub(session.ExpireAt) > 0 {
			http.Redirect(w, req, fmt.Sprintf("/auth/start?p=%s", path), http.StatusFound)
			return
		}
		// else ... what's happen
	}

	// render HTML (TODO template/html)
	t := `
<!doctype html>
<html>
<head><title></title></head>
<body>
<h2 style="text-align:center;">＼ You need to login ／</h2>
<form method="get" action="/auth/start">
<input type="hidden" name="p" value="` + path + `" />
<p style="text-align:center">
<input type="submit" value="" style="background:url(data:image/png;base64,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); width:191px; height:46px; cursor: pointer; border:0px none;" />
</p>
</form>
</body>
</html>
`

	io.WriteString(w, t)

}

type AuthSession struct {
	CSRFToken string
	Path      string
}

func (as *AuthSession) String() string {
	b := strings.Builder{}
	b.WriteString(as.CSRFToken)
	b.WriteString("|")
	b.WriteString(as.Path)

	return b.String()
}

func (as *AuthSession) Restore(data string) error {
	ss := strings.Split(data, "|")
	as.CSRFToken = ss[0]
	as.Path = ss[1]
	if as.Path == "" {
		as.Path = "/"
	}

	return nil
}

type IDSession struct {
	Issuer  string
	Subject string
	Email   string
}

func (is *IDSession) String() string {
	b := strings.Builder{}
	b.WriteString(is.Issuer)
	b.WriteString("('-'o)")
	b.WriteString(is.Subject)
	b.WriteString("('-'o)")
	b.WriteString(is.Email)

	return b.String()
}

func (is *IDSession) Restore(data string) error {
	ss := strings.Split(data, "('-'o)")
	is.Issuer = ss[0]
	is.Subject = ss[1]
	is.Email = ss[2]

	return nil
}

func escapeRequestPath(u *url.URL) string {
	var b strings.Builder
	b.WriteString(url.PathEscape(u.RequestURI()))
	if u.RawQuery != "" {
		b.WriteString("%3F")
		b.WriteString(url.QueryEscape(u.RawQuery))
	}

	return b.String()
}
