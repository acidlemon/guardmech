package generic

import (
	"context"

	"github.com/acidlemon/guardmech/backend/oidconnect"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type provider struct {
	p *oidc.Provider
}

func NewProvider(ctx context.Context, issuer string) (oidconnect.OIDCProvider, error) {
	p, err := oidc.NewProvider(ctx, issuer)
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
