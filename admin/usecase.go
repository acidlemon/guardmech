package admin

import (
	"database/sql"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
)

type AdminService interface {
	CreatePrincipal(Context, *db.Tx, string, string) (*membership.Principal, error)
	CreateAPIKey(Context, *db.Tx, *membership.Principal, string) (*membership.APIKey, string, error)

	FindPrincipalBySeqID(Context, *sql.Conn, int64) (*membership.Principal, error)

	FetchAllPrincipal(Context, *sql.Conn) ([]*membership.Principal, error)
	FetchAllRole(Context, *sql.Conn) ([]*membership.Role, error)
}

type Usecase struct {
	repos membership.Repository
	svc   AdminService
}

func (u *Usecase) CreatePrincipal(ctx Context, name, description string) (*membership.Principal, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	pri, err := u.svc.CreatePrincipal(ctx, tx, name, description)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return pri, err
}

func (u *Usecase) ShowPrincipal(ctx Context, ID int64) (*membership.PrincipalPayload, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	payload, err := u.repos.FetchPrincipalPayload(ctx, conn, ID)
	return payload, err
}

func (u *Usecase) ListPrincipals(ctx Context) ([]*membership.Principal, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	list, err := u.repos.FetchAllPrincipal(ctx, conn)

	return list, err
}

func (u *Usecase) ListRoles(ctx Context) ([]*membership.Role, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	list, err := u.repos.FetchAllRole(ctx, conn)

	return list, err
}

func (u *Usecase) CreateAPIKey(ctx Context, principalID int64, name string) (*membership.APIKey, string, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, "", err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	pri, err := u.svc.FindPrincipalBySeqID(ctx, conn, principalID)
	if err != nil {
		return nil, "", err
	}

	ap, rawToken, err := u.svc.CreateAPIKey(ctx, tx, pri, name)
	if err != nil {
		return nil, "", err
	}

	tx.Commit()

	return ap, rawToken, err
}
