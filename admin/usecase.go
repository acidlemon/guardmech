package admin

import (
	"context"
	"database/sql"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
)

type AdminService interface {
	CreatePrincipal(context.Context, *db.Tx, string, string) (*membership.Principal, error)
	CreateAPIKey(context.Context, *db.Tx, *membership.Principal, string, string) (*membership.APIKey, error)

	FetchAllPrincipal(context.Context, *sql.Conn) ([]*membership.Principal, error)
	FetchAllRole(context.Context, *sql.Conn) ([]*membership.Role, error)
}

type Usecase struct {
	repos membership.Repository
	svc   AdminService
}

func (u *Usecase) CreatePrincipal(ctx context.Context, name, description string) (*membership.Principal, error) {
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

func (u *Usecase) ShowPrincipal(ctx context.Context, ID int64) (*membership.PrincipalPayload, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	payload, err := u.repos.FetchPrincipalPayload(ctx, conn, ID)
	return payload, err
}

func (u *Usecase) ListPrincipals(ctx context.Context) ([]*membership.Principal, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	list, err := u.repos.FetchAllPrincipal(ctx, conn)

	return list, err
}

func (u *Usecase) ListRoles(ctx context.Context) ([]*membership.Role, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	list, err := u.repos.FetchAllRole(ctx, conn)

	return list, err
}

func (u *Usecase) CreateAPIKey(ctx context.Context, principalID int64, name, description string) (*membership.APIKey, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	pri, err := u.svc.

	ap, err := u.svc.CreateAPIKey(ctx, tx, principalID, name, description)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return ap, err
}
