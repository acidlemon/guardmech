package membership

import (
	"log"

	"github.com/acidlemon/guardmech/app/logic/auth"
	"github.com/google/uuid"
)

type Service struct {
	q Query
	//cmd Command
}

func NewService(q Query) *Service {
	return &Service{
		q: q,
	}
}

func (s *Service) FindPrincipalByID(ctx Context, principalID string) (*Principal, error) {
	pri, err := s.q.FindPrincipals(ctx, []string{principalID})
	if err != nil {
		return nil, err
	}
	if len(pri) == 0 {
		return nil, ErrNoEntry
	}

	return pri[0], nil
}

func (s *Service) FindPrincipalByOIDC(ctx Context, issuer, subject string) (*Principal, error) {
	pri, err := s.q.FindPrincipalIDByOIDC(ctx, issuer, subject)
	if err != nil {
		return nil, err
	}

	return pri, nil
}

func (s *Service) EnumeratePrincipalIDs(ctx Context) ([]string, error) {
	IDs, err := s.q.EnumeratePrincipalIDs(ctx)
	if err != nil {
		return nil, err
	}

	log.Println(IDs)

	if len(IDs) == 0 {
		return []string{}, nil
	}

	result := make([]string, 0, len(IDs))
	for _, u := range IDs {
		result = append(result, u.String())
	}

	return result, nil
}

func (s *Service) CreatePrincipalFromOpenID(ctx Context, token *auth.OpenIDToken) (*Principal, *OIDCAuthorization, error) {
	a := &OIDCAuthorization{
		OIDCAuthID: uuid.New(),
		Issuer:     token.Issuer,
		Subject:    token.Sub,
		Email:      token.Email,
		Name:       token.Name,
	}
	pri := &Principal{
		PrincipalID: uuid.New(),
		Name:        token.Name,
		Description: token.Email,
		auth:        a,
	}

	return pri, a, nil
}

func (s *Service) CreatePrincipalFromAPIKey(ctx Context, apiKey string) (*Principal, error) {
	//	s.cmd.CreatePrincipal()
	return &Principal{}, nil
}

func (s *Service) SetupPrincipalAsOwner(ctx Context, pri *Principal) (*Role, *Permission, error) {

	r, err := pri.AttachNewRole(RoleOwner, "")
	if err != nil {
		return nil, nil, err
	}

	perm, err := r.AttachNewPermission(ctx, PermissionOwner, "")
	if err != nil {
		return nil, nil, err
	}

	return r, perm, nil

}
