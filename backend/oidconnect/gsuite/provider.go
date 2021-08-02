package gsuite

import (
	"context"

	"github.com/acidlemon/guardmech/backend/oidconnect"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type provider struct {
	p *oidc.Provider
}

func NewProvider(ctx context.Context) (oidconnect.OIDCProvider, error) {
	p, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	return &provider{
		p: p,
	}, nil
}

func (p *provider) Verifier(config *oidc.Config) *oidc.IDTokenVerifier {
	return p.p.Verifier(config)
}

func (p *provider) Endpoint() oauth2.Endpoint {
	return p.p.Endpoint()
}
