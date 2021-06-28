package gsuite

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Provider struct {
	p *oidc.Provider
}

func New(ctx context.Context) (*Provider, error) {
	p, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	return &Provider{
		p: p,
	}, nil
}

func (p *Provider) Verifier(config *oidc.Config) *oidc.IDTokenVerifier {
	return p.p.Verifier(config)
}

func (p *Provider) Endpoint() oauth2.Endpoint {
	return p.p.Endpoint()
}
