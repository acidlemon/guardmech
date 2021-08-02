package oidconnect

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCProvider interface {
	Verifier(config *oidc.Config) *oidc.IDTokenVerifier
	Endpoint() oauth2.Endpoint
}

type GroupInquirer interface {
	IsMember(ctx context.Context, email, group string) (bool, error)
}
