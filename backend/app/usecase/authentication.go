package usecase

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/acidlemon/guardmech/backend/app/config"
	"github.com/acidlemon/guardmech/backend/app/logic/auth"
	"github.com/acidlemon/guardmech/backend/app/logic/membership"
	"github.com/acidlemon/guardmech/backend/app/usecase/payload"
	"github.com/acidlemon/guardmech/backend/db"
	"github.com/acidlemon/guardmech/backend/oidconnect"
	"github.com/acidlemon/guardmech/backend/oidconnect/gsuite"
	"github.com/acidlemon/guardmech/backend/persistence"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var clientID string
var clientSecret string
var redirectURL string

func init() {
	clientID = os.Getenv("GUARDMECH_CLIENT_ID")
	clientSecret = os.Getenv("GUARDMECH_CLIENT_SECRET")
	redirectURL = os.Getenv("GUARDMECH_REDIRECT_URL")
}

type Authentication struct {
	conf     *oauth2.Config
	provider oidconnect.OIDCProvider
	inquirer oidconnect.GroupInquirer
}

func NewAuthentication() *Authentication {
	ctx := context.Background()
	//	p, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	p, err := gsuite.NewProvider(ctx)
	if err != nil {
		// handle error
		panic(`failed to initialize Google OpenID Connect provider:` + err.Error())
	}

	// use GOOGLE_APPLICATION_CREDENTIALS
	q, err := gsuite.NewGroupInquirer(ctx)
	if err != nil {
		// handle error
		//panic(`failed to initialize Google Groups Inquirer:` + err.Error())
		q = nil
	}

	return &Authentication{
		conf: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     p.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		},
		provider: p,
		inquirer: q,
	}
}

func (u *Authentication) StartAuth(returnPath string) (*AuthSession, time.Time, string) {
	a := auth.NewAuthenticator(u.conf, u.provider)
	state, url, expireAt := a.StartAuthentication()

	as := &AuthSession{
		CSRFToken: state,
		Path:      returnPath,
	}

	return as, expireAt, url
}

func (u *Authentication) VerifyAuth(ctx Context, as *AuthSession, state, code string) (is *IDSession, expireAt time.Time, path string, reserr error) {
	// decode path
	path, err := url.PathUnescape(as.Path)
	if err != nil {
		path = "/"
	}

	// CSRF check
	if state != as.CSRFToken {
		reserr = securityError("Detect Possibility of CSRF Attack", nil)
		return
	}

	a := auth.NewAuthenticator(u.conf, u.provider)
	token, err := a.VerifyAuthentication(ctx, code)
	if err != nil {
		reserr = securityError("Failed to Verify AuthCode", err)
		return
	}

	// if first user, set as owner
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		reserr = systemError("Could not start transaction", err)
		return
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	manager := membership.NewManager(q)

	// generate membership token
	pri, err := manager.FindPrincipalByOIDC(ctx, token.Issuer, token.Sub)
	if err != nil {
		list, err := manager.EnumeratePrincipalIDs(ctx)
		if err != nil {
			reserr = systemError("Failed to Check Principal", err)
			return
		}

		cmd := persistence.NewCommand(tx)

		if len(list) == 0 {
			// 最初のユーザー
			log.Println("No User!! Entering setup mode.")

			var oidc *membership.OIDCAuthorization
			pri, oidc, err = manager.CreatePrincipalFromOpenID(ctx, token)
			if err != nil {
				log.Println("failed to setup:", err.Error())
				reserr = systemError("Failed to Setup", err)
				return
			}

			g, r, perm, err := manager.SetupPrincipalAsOwner(ctx, pri)
			if err != nil {
				reserr = systemError("Failed to Setup", err)
				return
			}

			roperm, err := manager.SetupSystemMembership(ctx)
			if err != nil {
				reserr = systemError("Failed to Setup", err)
				return
			}

			cmd.SavePermission(ctx, roperm)
			cmd.SavePermission(ctx, perm)
			cmd.SaveRole(ctx, r)
			cmd.SaveGroup(ctx, g)
			cmd.SavePrincipal(ctx, pri)
			cmd.SaveAuthOIDC(ctx, oidc, pri)
			if cmd.Error() != nil {
				reserr = systemError("Failed to save item", err)
				return
			}

		} else {
			// TODO ここで2人目以降のユーザを追加するための処理が必要
			rules, err := manager.EnumerateMappingRules(ctx)
			if err != nil {
				reserr = systemError("Failed to enumerate mapping rules", err)
				return
			}

			gi := auth.NewGroupInquirer(u.inquirer)
			ruleman := membership.NewMappingRuleManager(rules, gi)
			foundRules, err := ruleman.FindMatchedRules(ctx, token)
			if err != nil {
				reserr = verificationError("You have no matching rules.", err)
				return
			}

			var oidc *membership.OIDCAuthorization
			pri, oidc, err = manager.CreatePrincipalFromRules(ctx, token, foundRules)
			if err != nil {
				log.Println("failed to prepare principal:", err.Error())
				reserr = systemError("Failed to create principal", err)
				return
			}

			cmd.SavePrincipal(ctx, pri)
			cmd.SaveAuthOIDC(ctx, oidc, pri)
			if err = cmd.Error(); err != nil {
				log.Println("failed to save principal:", err)
				reserr = systemError("Failed to save item", err)
				return
			}

		}
	}

	err = tx.Commit()
	if err != nil {
		reserr = systemError("Failed to commit transaction", err)
		return
	}

	p := payload.SessionPrincipalFromEntity(pri)
	now := time.Now()
	is = &IDSession{
		Issuer:  token.Issuer,
		Subject: token.Sub,
		Email:   token.Email,
		Membership: MembershipToken{
			NextCheck: now.Add(1 * time.Minute), // TODO 使わなそう
			Principal: p,
		},
	}

	// OK!! Create Session!!
	expireAt = now.Add(config.SessionLifeTime)

	return
}

func (u *Authentication) Authorization(ctx Context, is *IDSession) (string, *payload.SessionPrincipal, error) {

	return is.Email, is.Membership.Principal, nil
}

func (u *Authentication) NeedAuthPrompt(ctx Context, expireAt time.Time) bool {
	if time.Now().Sub(expireAt) > 0 {
		return false
	}

	return true
}

func (u *Authentication) verifyAuthCode(ctx Context, code string) (*auth.OpenIDToken, error) {
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
	var claims auth.OpenIDToken
	if err := idToken.Claims(&claims); err != nil {
		// handle error
		return nil, err
	}

	return &claims, nil
}
