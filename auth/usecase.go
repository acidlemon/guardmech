package auth

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
	oidc "github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var randLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateRandomString(length int, letters []rune) string {
	b := strings.Builder{}
	b.Grow(length)
	lc := len(letters)

	for i := 0; i < length; i++ {
		b.WriteRune(letters[rand.Intn(lc)])
	}
	return b.String()
}

type Usecase struct {
	conf     *oauth2.Config
	provider *oidc.Provider
}

func (u *Usecase) StartAuth(returnPath string) (*SessionPayload, string) {
	csrfToken := generateRandomString(32, randLetters)
	as := &AuthSession{
		CSRFToken: csrfToken,
		Path:      returnPath,
	}

	authSession := &SessionPayload{
		ExpireAt: time.Now().Add(5 * time.Minute),
		Data:     as,
	}

	url := u.conf.AuthCodeURL(csrfToken)

	return authSession, url
}

func (u *Usecase) CallbackAuth(ctx context.Context, cookieValue, state, code string) (session *SessionPayload, path string, reserr error) {
	var svc AuthService
	svc = membership.NewService()

	as := &AuthSession{}
	_, err := RestoreSessionPayload(cookieValue, as)
	if err != nil {
		log.Println(err)
		reserr = NewHttpError(http.StatusForbidden, "Session Validation Failed", err)
		return
	}

	// decode path
	path, err = url.PathUnescape(as.Path)
	if err != nil {
		path = "/"
	}

	// CSRF check
	if state != as.CSRFToken {
		reserr = NewHttpError(http.StatusForbidden, "Detect Possibility of CSRF Attack", nil)
		return
	}

	// ID Token
	token, err := u.verifyAuthCode(ctx, code)
	if err != nil {
		reserr = NewHttpError(http.StatusForbidden, "Failed to Verify AuthCode", err)
		return
	}

	// normalize issuer (remove "https://" for Google)
	token.Issuer = strings.Replace(token.Issuer, "https://", "", -1)

	// if first user, set as owner
	conn, err := db.GetConn(ctx)
	if err != nil {
		reserr = NewHttpError(http.StatusInternalServerError, "Could not GetConn", err)
		return
	}
	defer conn.Close()

	ok, err := svc.HasPrincipal(ctx, conn)
	if err != nil {
		reserr = NewHttpError(http.StatusInternalServerError, "Failed to Check Principal", err)
		return
	}
	// first cut
	if !ok {
		log.Println("No User!! Entering setup mode.")

		err := svc.CreateFirstPrincipal(ctx, conn, token)
		if err != nil {
			log.Println("failed to setup:", err.Error())
			reserr = NewHttpError(http.StatusInternalServerError, "Failed to Setup", err)
			return
		}
	}

	// generate membership token
	pri, err := svc.FindPrincipal(ctx, conn, token.Issuer, token.Sub)
	if err != nil {
		reserr = NewHttpError(http.StatusForbidden, "Could Not Find Principal", err)
		return
	}

	payload, err := svc.FetchPrincipalPayload(ctx, conn, pri.ID)
	if err != nil {
		reserr = NewHttpError(http.StatusForbidden, "Could Not Find PrincipalPayload", err)
		return
	}

	now := time.Now()
	is := &IDSession{
		Issuer:  token.Issuer,
		Subject: token.Sub,
		Email:   token.Email,
		Membership: MembershipToken{
			NextCheck: now.Add(5 * time.Minute), // TODO 使わなそう
			Principal: payload,
		},
	}

	// OK!! Create Session!!
	session = &SessionPayload{
		Data:     is,
		ExpireAt: now.Add(5 * time.Minute),
	}
	return
}

func (u *Usecase) Authorization(ctx context.Context, cookieValue string) (string, *membership.PrincipalPayload, error) {
	is := &IDSession{}
	session, err := RestoreSessionPayload(cookieValue, is)
	if err != nil {
		return "", nil, NewHttpError(http.StatusUnauthorized, "failed to restore cookie", err)
	}

	// check ExpireAt
	if time.Now().Sub(session.ExpireAt) > 0 {
		return "", nil, NewHttpError(http.StatusUnauthorized, "session expired", nil)
	}

	return is.Email, is.Membership.Principal, nil
}

func (u *Usecase) NeedAuthPrompt(ctx context.Context, cookieValue string) bool {
	// cookie check
	is := &IDSession{}
	session, err := RestoreSessionPayload(cookieValue, is)
	if err == nil {
		// cookie is live, try to authenticate without prompt
		if time.Now().Sub(session.ExpireAt) > 0 {
			return false
		}
	}

	return true
}

func (u *Usecase) verifyAuthCode(ctx context.Context, code string) (*membership.OpenIDToken, error) {
	var verifier = u.provider.Verifier(&oidc.Config{ClientID: u.conf.ClientID})
	oauth2Token, err := u.conf.Exchange(ctx, code)
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
	var claims membership.OpenIDToken
	if err := idToken.Claims(&claims); err != nil {
		// handle error
		return nil, err
	}

	return &claims, nil
}
