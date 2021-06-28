package oidconnect

import (
	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type Membership interface {
	Verifier(config *oidc.Config) *oidc.IDTokenVerifier
	Endpoint() oauth2.Endpoint
}
