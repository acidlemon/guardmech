package membership

import (
	"context"
	"database/sql"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) HasPrincipal(ctx context.Context, conn *sql.Conn) (bool, error) {
	return HasPrincipal(ctx, conn)
}
func (s *Service) CreateFirstPrincipal(ctx context.Context, conn *sql.Conn, idToken *OpenIDToken) error {
	return CreateFirstPrincipal(ctx, conn, idToken)
}
func (s *Service) FindPrincipal(ctx context.Context, conn *sql.Conn, issuer, subject string) (*Principal, error) {
	return FindPrincipal(ctx, conn, issuer, subject)
}
func (s *Service) FetchPrincipalPayload(ctx context.Context, conn *sql.Conn, id int64) (*PrincipalPayload, error) {
	return FetchPrincipalPayload(ctx, conn, id)
}

func (s *Service) FetchAllPrincipal(ctx context.Context, conn *sql.Conn) ([]*Principal, error) {
	return FetchAllPrincipal(ctx, conn)
}
func (s *Service) FetchAllRole(ctx context.Context, conn *sql.Conn) ([]*Role, error) {
	return FetchAllRole(ctx, conn)
}
