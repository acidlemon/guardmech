package usecase

import (
	"database/sql"

	"github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/guardmech/db"
)

type AdminService interface {
	CreatePrincipal(Context, *db.Tx, string, string) (*membership.Principal, error)
	CreateAPIKey(Context, *db.Tx, *membership.Principal, string) (*membership.AuthAPIKey, string, error)

	FindPrincipalBySeqID(Context, *sql.Conn, int64) (*membership.Principal, error)

	FetchAllPrincipal(Context, *sql.Conn) ([]*membership.Principal, error)
	FetchAllRole(Context, *sql.Conn) ([]*membership.Role, error)
}

type Administration struct {
	repos membership.Service
	svc   AdminService
}

func NewAdministration() *Administration {
	return &Administration{}
}

func (u *Administration) CreatePrincipal(ctx Context, name, description string) (*membership.Principal, error) {
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

/*
func (u *Administration) ShowPrincipal(ctx Context, ID int64) (*membership.Principal, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	payload, err := u.repos.FetchPrincipalPayload(ctx, conn, ID)
	return payload, err
}

func (u *Administration) ListPrincipals(ctx Context) ([]*membership.Principal, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	list, err := u.repos.FetchAllPrincipal(ctx, conn)

	return list, err
}

func (u *Administration) ListRoles(ctx Context) ([]*membership.Role, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	list, err := u.repos.FetchAllRole(ctx, conn)

	return list, err
}
*/

func (u *Administration) CreateAPIKey(ctx Context, principalID int64, name string) (*membership.AuthAPIKey, string, error) {
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

func (u *Administration) CreateRole(ctx Context) (*membership.Role, error) {
	return nil, nil
}
func (u *Administration) FetchRole(ctx Context, id int64) (*membership.Role, error) {
	return nil, nil
}
func (u *Administration) UpdateRole(ctx Context, id int64) (*membership.Role, error) {
	return nil, nil
}
func (u *Administration) ListMappingRules(ctx Context) ([]*membership.MappingRule, error) {
	return nil, nil
}
func (u *Administration) CreateMappingRule(ctx Context) (*membership.MappingRule, error) {
	return nil, nil
}
func (u *Administration) FetchMappingRule(ctx Context, id int64) (*membership.MappingRule, error) {
	return nil, nil
}
func (u *Administration) UpdateMappingRule(ctx Context, id int64) (*membership.MappingRule, error) {
	return nil, nil
}
func (u *Administration) ListGroups(ctx Context) ([]*membership.Group, error) {
	return nil, nil
}
func (u *Administration) CreateGroup(ctx Context) (*membership.Group, error) {
	return nil, nil
}
func (u *Administration) FetchGroup(ctx Context) (*membership.Group, error) {
	return nil, nil
}
func (u *Administration) UpdateGroup(ctx Context) (*membership.Group, error) {
	return nil, nil
}
func (u *Administration) ListPermissions(ctx Context) ([]*membership.Permission, error) {
	return nil, nil
}
func (u *Administration) CreatePermission(ctx Context) (*membership.Permission, error) {
	return nil, nil
}
func (u *Administration) FetchPermission(ctx Context) (*membership.Permission, error) {
	return nil, nil
}
func (u *Administration) UpdatePermission(ctx Context) (*membership.Permission, error) {
	return nil, nil
}
