package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/acidlemon/guardmech/backend/app/config"
	"github.com/acidlemon/guardmech/backend/app/usecase"
	"github.com/gorilla/mux"
)

const authSessionKey string = `_guardmech_csrf`
const sessionKey string = `_guardmech_session`

type AuthMux struct {
	u *usecase.Authentication
}

var clientID string
var clientSecret string
var redirectURL string

func init() {
	clientID = os.Getenv("GUARDMECH_CLIENT_ID")
	clientSecret = os.Getenv("GUARDMECH_CLIENT_SECRET")
	redirectURL = os.Getenv("GUARDMECH_REDIRECT_URL")
}

func NewAuthMux() *AuthMux {
	u := usecase.NewAuthentication()
	am := &AuthMux{
		u: u,
	}
	return am
}

func (a *AuthMux) Mux() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/auth/start", a.StartAuth)
	m.HandleFunc("/auth/callback", a.CallbackAuth)
	m.HandleFunc("/auth/request", a.AuthRequest)
	m.HandleFunc("/auth/sign_in", a.SignIn)
	m.HandleFunc("/auth/unauthorized", a.Unauthorized)

	return m
}

func (a *AuthMux) StartAuth(w http.ResponseWriter, req *http.Request) {
	log.Println("start auth")
	// return path
	if err := req.ParseForm(); err != nil {
		WriteHttpError(w, NewHttpError(http.StatusBadRequest, "Body Parsing Error", err))
		return
	}
	path := req.Form.Get("p")
	if path == "" {
		path = "%2F"
	}

	as, expireAt, url := a.u.StartAuth(path)
	session := &SessionPayload{
		Data:     as,
		ExpireAt: expireAt,
	}

	// bake CSRF session cookie
	http.SetCookie(w, session.MakeCookie(req, authSessionKey, 0))
	http.Redirect(w, req, url, http.StatusFound)
}

func (a *AuthMux) CallbackAuth(w http.ResponseWriter, req *http.Request) {
	log.Println("callback auth")
	// CSRF session validation
	c, err := req.Cookie(authSessionKey)
	if err != nil {
		WriteHttpError(w, NewHttpError(http.StatusForbidden, "No CSRF Session", err))
		return
	}

	req.ParseForm()
	state := req.Form.Get("state")
	code := req.URL.Query().Get("code")

	as := &usecase.AuthSession{}
	_, err = RestoreSessionPayload(c.Value, as)
	if err != nil {
		log.Println(err)
		WriteHttpError(w, fmt.Errorf("Session Validation Failed: %s", err))
		return
	}

	is, expireAt, path, err := a.u.VerifyAuth(req.Context(), as, state, code)

	// delete CSRF session cookie
	domain := req.URL.Host
	http.SetCookie(w, revokeCookie(domain, authSessionKey))

	if err != nil {
		WriteHttpError(w, NewHttpErrorFromErr(err))
		return
	}

	session := &SessionPayload{
		Data:     is,
		ExpireAt: expireAt,
	}

	http.SetCookie(w, session.MakeCookie(req, sessionKey, config.CookieLifeTime))
	http.Redirect(w, req, path, http.StatusFound)
}

func (a *AuthMux) AuthRequest(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie(sessionKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	is := &usecase.IDSession{}
	session, err := RestoreSessionPayload(c.Value, is)
	if err != nil {
		httperr := NewHttpError(http.StatusUnauthorized, "failed to restore cookie", err)
		WriteHttpError(w, httperr)
		return
	}

	// check ExpireAt
	if time.Now().Sub(session.ExpireAt) > 0 {
		log.Println("session expired")
		httperr := NewHttpError(http.StatusUnauthorized, "session expired", nil)
		WriteHttpError(w, httperr)
		return
	}

	email, principal, err := a.u.Authorization(req.Context(), is)
	if err != nil {
		httperr := NewHttpError(http.StatusUnauthorized, err.Error(), nil)
		WriteHttpError(w, httperr)
		return
	}

	// OK! print headers
	w.Header().Set("X-Guardmech-Email", email)
	if len(principal.Groups) > 0 {
		groups := make([]string, 0, len(principal.Groups))
		for _, v := range principal.Groups {
			groups = append(groups, v)
		}
		w.Header().Set("X-Guardmech-Groups", strings.Join(groups, ";"))
	}
	if len(principal.Roles) > 0 {
		roles := make([]string, 0, len(principal.Roles))
		for _, v := range principal.Roles {
			roles = append(roles, v)
		}
		w.Header().Set("X-Guardmech-Roles", strings.Join(roles, ";"))
	}
	if len(principal.Permissions) > 0 {
		perms := make([]string, 0, len(principal.Permissions))
		for _, v := range principal.Permissions {
			perms = append(perms, v)
		}
		w.Header().Set("X-Guardmech-Permissions", strings.Join(perms, ";"))
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *AuthMux) SignIn(w http.ResponseWriter, req *http.Request) {
	// catch path
	originUri := req.Header.Get("X-Auth-Request-Redirect")
	path := url.PathEscape(originUri)

	c, err := req.Cookie(sessionKey)
	if err == nil {
		is := &usecase.IDSession{}
		session, err := RestoreSessionPayload(c.Value, is)
		if err != nil {
			http.Redirect(w, req, fmt.Sprintf("/auth/start?p=%s", path), http.StatusFound)
			return
		}

		necessary := a.u.NeedAuthPrompt(req.Context(), session.ExpireAt)
		if !necessary {
			http.Redirect(w, req, fmt.Sprintf("/auth/start?p=%s", path), http.StatusFound)
			return
		}
	}

	// render HTML (TODO template/html)
	t := `
<!doctype html>
<html>
<head>
  <title>Need Authentication</title>
  <link rel="stylesheet" href="https://newcss.net/new.min.css">
</head>
<body>
<header>
  <h1>Authorization Required</h1>
</header>
<main>
  <h2>You need to login</h2>
  <form method="get" action="/auth/start">
    <input type="hidden" name="p" value="` + path + `" />
    <p>
      <input type="submit" value="" style="background:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAL8AAAAuCAYAAAB50MjgAAAAAXNSR0IArs4c6QAAD0lJREFUeAHtXQt0VNW5/mbmzCOTCUQSQoJBQ0NijAHUAEJraAERWUixLNcS7L2XdhWsVWutfVBN24VrVS3UltZKl66mVaxU6b2ktyK9lEa4baQqLwuFVAIRkTSYGCE0k3memdNvnzOPM5MJTCaJhPT8a53M2Xv/Zz++/e3///eeMwAYYiBgIGAgYCDw74WAKWm4Ih29koqMpIHAJY2Awt5HL3UgevKbZn3r6G327MINZrOl6JIeptF5A4EkBMLh0Olgd/t9u9eV/5ZFYhFA0ulY7M5xP/v+MmfhDZP02TqNi3T75nEZD//Gf5FaN5odCQgIg251FWzgWF7mJYsxmcWfiEhmizTsiC/6NtwWYxQw4/PSQkDwmz2OWXY9+fUh0KU1KqO3BgLpIxDjuZ786T9uaBoIjAAEDPKPgEk0hpAZAgb5M8PNeGoEIGCQfwRMojGEzBAwyJ8ZbsZTIwABg/wjYBKNIWSGgEH+zHAznhoBCMQO/AdjLMGjTfBtrYd84jhCba1QggFY8sbCOvla2D81H7ZpMwejGaMOA4FBQWBQyB9qfx/d6x5B8OD+Xp0K/eMUxOXbvhVS1VSMWr0GlqLLe+kZGQYCHzUCAw57gk1/w9l7P5eS+MmDkQ8fRPf6x5OzjbSBwEVBYECWP9TZgXMPPwDF3Z3QeUvxFZBKy6GEw5BbmhFmCCREKqvAqO88lqA7lImi8RYsrbQgl6PsORdGw34ZhwNai6UVVjww3Yx9r/qxsW2Qe2Ezo/YzVjg7ZNTuCg1y5YnVLZ1jx5wxYTyzJYjDiUVqqqjKhuc+bUXXUT/u2KK+z5VCK0VWrhmrpkkoIHZBfxgH3pLR0JVCb6iysiVs/ooduR8Gcc8zAbQMQTsDIn/3ujUJxDePLYDr/tWwz6pJ6Kq/cRfDnpeR861HYM4ZlVA2VInlS7Kw8ppEx7Z4jg1/2OzFuhYFN063oXICUD4vhI2/6gcp0uhw0cesmFtKaEtNqCL5U5EyjWrSUDHh5lkSyqj5yVySv8uEVUtsqM5WsOl/AmjkQndaTeqbXA57GtVFVKbPsuN7c6T4G2DMX8C8ZQd8WLl9aBezvpcOJqQcE7L0mYN4nzH5lbOvwVpYj6ClGAiZYc7LR+4Pn4FlfO943l4zB+L6qCSn1BYhvoKmAwH87hhw42wbaopMWHC7Db9f68cOWvyJtPx7/zK4xBdjPP12AC8dUJDdTk8zpIMmyf8YwLzRCrZHrHJlGReDTcHHrBr5Y82ny9mxVqyJEL/z3SCefzOEsmutWHyVBROvt2PFGx5s/Ig8gDoz7Lc3NojBvcmY/OH3X4B96hlYCr3oqZ+I7C/cm5L4g9vd9GpzWSN63hAe3i5DBGUNtPabVjtQaDGjOhd4p9SCq680I9hiwrYPxG8bTKj9TwdmTzBDCilo7VDgyjXBfcyPFbuAui/YcZk/hLdpXafxWVDn1JEAvvqKVn9Cz7LZBi1/zhgFRW+FgQobnloowdcVQrfNgrI8E+RAGK9u92HdYfV3FfHHx0rY+B82ONqDWPnroKr/07vtKOiMp59musgj48G6IGZXSJg8SsHUgwoeWm7HBJuoyoQ7v+TExO1e/CpSsyNfwtNftGlte8P4v5d9+DExSZalH5cgLK7cEcAdbF+VlhBsK7OxoMCEWVebsPF1Da9Vt9pxa4UFLrbpc4fwxx1+/PjteJ1V19nwjRoJhS6+SEm8mg8GYvMh6l3EkO2uGRJchLPzlIyzTguKQmJcqcOcqiorvjHXisIsrb6D+/345gDCyozJr5zbo+IijfMiZ1U7rPNuUdP6P4+/fP4foNxYbkENJ2+w5fQ/FfXXClKWhOc+T9LvDKL+ZAifXdsTa2rFBAvyOSmFEZ/6wEon5hZoxe6ACcX0EkIco7TPsdR1uSTMzAPcNEUuTsDEKXZ87e8y1rRoz8X/mlA0mvpZZlAdXrruXDFh7E8hSUCekDBmLFhEL3TYn+gdzimQqJtfYsENCOLNcgmVgjxsezbT29jvMpEmtB2sOz/fzLoV5BFGK0kURVPifdQGiH5Jo/kc2/axbQf7tfh2O/av9aFRFOokL1sbb/PBRI+4rq4H63R6K+7IwjKGdULcXhoKMnjxUidy6z1YwwVQxAX/k4VaD3wsF2OqpOd4YTSwZLOMqul2PMiQTYgoz58gIV8kiH2qMEd48x9y7yKekAOsz2ZC9SwHnvZ5cffrNDAZSGJQ3J8KfO/FtM15V8FkIdpJ0nCYm6TzXH8lIYdE2oJ4co9Wd26RFfd+1olXVzux/iYLciINBvRNc3NVoxJfQf2mHixZ34O6IxqgUQponwq2Rsq3CuZR8sZoBNBS8b8JLjtaiVvGV9Z6sISkaxXtE7NZ9EIJwo4d/FDkWFA1Fph/TRRXMz5ZAUy/Qpuy1mbN40RsMwL0civWe7FfjREUPP+UB7V6rxJpe1GsbTMqkttmq5c5E3oDYW3X3GpD7S029XPFlRwv8VqqEj+KhyeGV808KzE24b6bNeK3cp+waL0Ht28KwM2qXSTxci6w5RHiv39IK7+tPqj9vIq4pApz7r9JI/6JPV4seMIDoS+kbKYVpepd//9kTv7+t9XriUCUFL1KBp6xrcGHeT/z4Q9HQ+gSJzwWE6bMcOCFO6O2Md5GTj6ttEiSQPUntfz/b0lhTVi+OVJ+1h937/Gazn/nbo9ufsNoj5w6RT4SHnydYYaQiiouyvHxxVV6jRWfKBZTpuCtoyn6x5KotbdFb0RFlPTaNmG0U2tPjsxNNUOXmincwF9vVT9njQdyCs1qaIQPZfwigseLDUGV3KCFv5LtCS8k+rmzURtL90kZ70QGO47WP1q+IxK2dJ8Ka8+Lx3oJPWFkUU6Y7MDm+7NQF/EqtBEpPUWvKlJk9GZCCqWUWY4JPD88ohb5PCeQFQ7BYo5aKe2J4iSrGCQO7XTrUcmNuNhoerA+qxhK3cSw5djhINZt0cCv4WnFt7mRc5XYsCJbRirSCSCjhs+ZRJ5o36LlgZQVRLX6+EyEpw8loLEpBN8Mhinss5DWPX6cKLWj5iobFoqMUAi7+ziejXoCoZYgabWt4HC7gpkMq8aME4sshB1bPWgiFp+4JQuLJ5gQ5KIoYCiYijh6WybmWuw9XAKweLSp5mUzromV94GzeDqlcBwS65W4sDrP8Z5GKJWnSPlsUmaqMSSppE6Gc6bBTPI3y6Pw0JlK3HWyEYsmfipBeePdUapo2f+7L4if7oiz5oq8oXE81dfbsLiEm8oSYNuzGh0aD4XgJvmFpxdWMd4LoLuTsTDzXfQODFXRwsLx3OxeNGkL4RT7UKZuXoG3GH/v5nTX5FlU0rm5Odx7gc65PRdQ6KO4tYMehRv6YmK4/A0vXuxScJq6lZoNUZ9qaSaW3MC7+F3ADcxp4FUlvk8RpYzHRdSmWXbgymLiKA4UuMcZoy5ABR1UKIssxkllLN/LcnocsdHuS6KLunm3F19WN9wMAa8jfY/IGX8HkDn5xy3H79/diR+4p5BIFjx1aBOmjr0axa5xKft/zqNgy97oELR/HGiGODUZAql/U8adJYwRi3iuf78FzZ0kMzeK2uSE8HoXMFnfbg/DEQakM7lpe5Au9ZY2BZUlur5xsWQEFKsQm7fzWaYIv/W94X0Y+9mHMi5gMDbc/QH5EfEGgiDvHYszMWo4o/VoaROWrXJg3HY/tibVrE9Gn9HnNe7yY/+1TlRzU7zyHifmvcu2cnjMyROqmBCvPSTwXC7Gh76ehUU6vJo4x6dplf97XwjVNRZUL3SirlxGjtjQEg+5I4iXuKCu/RvL6d2mzHdiY5kMB8tV8utgV9tTMYzXVznHiU2TZJyxW1DJ0yfUmHHbkwH1RC/WvzRvMja9jvy5aMj+nEp80Van7yzu3vldNLbt69X00bMn8MDmd9F2Nh7yzKuyxOLLXg8MMKObpvue+gBayTqJhBZEziWI7nMyNtQlnq5orlpB7S98OEQ3KiyU0Ne7cHDhcr0kyD97tLEk6CVoMEHrnUD8OGdFRKGK3gPpH98ZIXjXqZBm5dtkNKuVKdjXFMdRCx/Upvi4gm0klRAXT3cqxqi32p+021bwzQ1evHFa7Cl4olUiacTnSdGxQ348Jqw023n0WerwOFjgNYV4CeNw7IA/ZpX3Nvqw4UBIxXEij33zudLcJP73nufxLXUbuSd7KbJvKWYb6kkP8wUuCac9EQxFfT/aI6seupALRRBf5rFZ3YuZEV80pVvOyJr7qNvz6sPZIj8t6fB8iGXbH0R3MCGowxWuIpRfNpF7TDNOdrfh7bPvwBR28Oz6Llg91chmKPvzlVkYNzr9tTfvscQ20uoglXK4rygQ5jCooOU8VRTxdYf7ShT8slFGR48J/3VnFpbS8nYe8eGO3+mZk27LF1GPYy7lmFtoYQckPE4sFZErHXZLZLEn1yfwFYcFbpYLUvcSUQc3uMKAJODPRbNmsYTX/szjXO4DK6byW+X5/FaZnu6uJ/znCWVYH124h/N5+jzz2asfkYydteouRDUjGXnzaMUFzjz8ZHYtvv7aWpzxC7OpyXvu0xCXXhSzD96iJzHKfxsemb28X8TX19Pf++6+JiWhIh7NLbRhJk3OtKttcDP2z1VjAgUNf7rEiC/GxTG3JIwvwwTjd7H/OZ9cEF9RB8O2ZLnpZrGBN/OS0MWQM1d8d0E5sT94gb6zvmQ3nFx5mun0TW8fFU7OL8fG+d9HdcE1fWjEs0tyLscTiz+O6/TxdLz4It4x7HmWX/gwvvWR+C4Lwxwe4z3/nAc/HySgL+LghmXTDa/4UHdARidtsPjC0Mdvnd9gaLNyV+oj3KEYxIDCnuQONZ05ji3Hd+D4uffQ6n4f/lAAeY5cVOWVY27xDZjDy8xQKBPJNOzJpC3jmZGLwKCFPckQVY6ZhMoZk5KzjbSBwLBEIDMzPCyHYnTKQKB/CBjk7x9ehvYIQsAg/wiaTGMo/UPAIH//8DK0RxACBvlH0GQaQ+kfAgb5+4eXoT2CEDDIP4Im0xhK/xAwyN8/vAztEYSAnvxKSPZ9IP7zt+Emw7FPww0joz8XRkDwm1qxt/30L7aFO5u2PrZ60621FmtW7A3TC1dpaBgIDH8EQkFvZ8ffX3mUPY29PKR/t0d4AfGPDVzOS/x+Tl/GpCEGApcsAsLai39K5B+8xA/N1AWQTHCxAMRvaZLzmWWIgcAljYBYAOL99Jjlv6RHY3TeQGAgCPwLjsfs3HaOnC0AAAAASUVORK5CYII=); width:191px; height:46px; cursor: pointer; border:0px none;" />
    </p>
  </form>
</main>
</body>
</html>
`
	io.WriteString(w, t)
}

func (a *AuthMux) Unauthorized(w http.ResponseWriter, req *http.Request) {
	WriteHttpError(w, NewHttpError(http.StatusForbidden, "No Access Permission", nil))
}
