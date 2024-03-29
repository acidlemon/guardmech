package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/acidlemon/guardmech/backend/app/config"
	"github.com/acidlemon/guardmech/backend/app/logic"
	"github.com/acidlemon/guardmech/backend/oidconnect"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	oidcConf *oauth2.Config
	provider oidconnect.OIDCProvider
}

func NewAuthenticator(conf *oauth2.Config, provider oidconnect.OIDCProvider) *Authenticator {
	return &Authenticator{
		oidcConf: conf,
		provider: provider,
	}
}

// OpenID Connectを利用した認証の開始
func (a *Authenticator) StartAuthentication() (string, string, time.Time) {
	state := logic.GenerateRandomString(32)
	url := a.oidcConf.AuthCodeURL(state)
	expireAt := time.Now().Add(config.AuthenticationTimeout)

	return state, url, expireAt
}

// OpenID Connectの認証結果の検証
func (a *Authenticator) VerifyAuthentication(ctx context.Context, code string) (*OpenIDToken, error) {
	var verifier = a.provider.Verifier(&oidc.Config{ClientID: a.oidcConf.ClientID})
	oauth2Token, err := a.oidcConf.Exchange(ctx, code)
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

	// extract access token
	// accessToken := oauth2Token.AccessToken

	// normalize issuer (remove "https://" for Google)
	claims.Issuer = strings.Replace(claims.Issuer, "https://", "", -1)

	return &claims, nil
}
