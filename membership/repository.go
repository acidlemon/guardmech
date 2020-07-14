package membership

import (
	"context"
	"database/sql"

	"github.com/acidlemon/guardmech/db"
)

type Repository interface {
	HasPrincipal(ctx context.Context, conn *sql.Conn) (bool, error)

	SavePrincipal(ctx context.Context, tx *db.Tx, pri *Principal) (int64, error)
	SaveAuth(ctx context.Context, tx *db.Tx, a *Auth) (int64, error)
	SaveAPIKey(ctx context.Context, tx *db.Tx, a *APIKey) (int64, error)
	SavePermission(ctx context.Context, tx *db.Tx, a *Permission) (int64, error)
	SaveRole(ctx context.Context, tx *db.Tx, a *Role) (int64, error)

	FindPrincipal(ctx context.Context, conn *sql.Conn, issuer, subject string) (*Principal, error)
	FindPrincipalBySeqID(ctx context.Context, conn *sql.Conn, principalID int64) (*Principal, error)
	FetchPrincipalPayload(ctx context.Context, conn *sql.Conn, id int64) (*PrincipalPayload, error)
	FetchAllPrincipal(ctx context.Context, conn *sql.Conn) ([]*Principal, error)
	FetchAllRole(ctx context.Context, conn *sql.Conn) ([]*Role, error)
}
