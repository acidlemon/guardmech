package membership

import (
	"github.com/acidlemon/guardmech/app/logic/auth"
	"github.com/google/uuid"
)

type Manager struct {
	q Query
}

func NewManager(q Query) *Manager {
	return &Manager{
		q: q,
	}
}

func (s *Manager) FindPrincipalByID(ctx Context, principalID string) (*Principal, error) {
	pri, err := s.q.FindPrincipals(ctx, []string{principalID})
	if err != nil {
		return nil, err
	}
	if len(pri) == 0 {
		return nil, ErrNoEntry
	}

	return pri[0], nil
}

func (s *Manager) FindPrincipalByOIDC(ctx Context, issuer, subject string) (*Principal, error) {
	pri, err := s.q.FindPrincipalIDByOIDC(ctx, issuer, subject)
	if err != nil {
		return nil, err
	}

	return pri, nil
}

func (s *Manager) EnumeratePrincipalIDs(ctx Context) ([]string, error) {
	IDs, err := s.q.EnumeratePrincipalIDs(ctx)
	if err != nil {
		return nil, err
	}

	if len(IDs) == 0 {
		return []string{}, nil
	}

	result := make([]string, 0, len(IDs))
	for _, u := range IDs {
		result = append(result, u.String())
	}

	return result, nil
}

func (m *Manager) FindPrincipals(ctx Context, ids []string) ([]*Principal, error) {
	list, err := m.q.FindPrincipals(ctx, ids)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Manager) CreatePrincipalFromOpenID(ctx Context, token *auth.OpenIDToken) (*Principal, *OIDCAuthorization, error) {
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

func (s *Manager) CreatePrincipalFromAPIKey(ctx Context, apiKey string) (*Principal, error) {
	//	s.cmd.CreatePrincipal()
	return &Principal{}, nil
}

func (s *Manager) SetupPrincipalAsOwner(ctx Context, pri *Principal) (*Role, *Permission, error) {

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
