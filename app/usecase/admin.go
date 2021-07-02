package usecase

import (
	"log"

	"github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/persistence"
)

type Administration struct {
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

	// pri, err := u.svc.CreatePrincipal(ctx, tx, name, description)
	// if err != nil {
	// 	return nil, err
	// }

	tx.Commit()

	return nil, err
}

func (u *Administration) ShowPrincipal(ctx Context, principalID string) (*membership.Principal, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, systemError("Could not start transaction", err)
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	manager := membership.NewManager(q)

	pri, err := manager.FindPrincipalByID(ctx, principalID)
	if err != nil {
		return nil, err
	}

	return pri, err
}

func (u *Administration) ListPrincipals(ctx Context) ([]*membership.Principal, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, systemError("Could not start transaction", err)
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	manager := membership.NewManager(q)

	ids, err := manager.EnumeratePrincipalIDs(ctx)
	if err != nil {
		return nil, err
	}

	list, err := manager.FindPrincipals(ctx, ids)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (u *Administration) CreateAPIKey(ctx Context, principalID string, name string) (*membership.AuthAPIKey, string, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, "", systemError("Could not start transaction", err)
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)
	pri, err := manager.FindPrincipalByID(ctx, principalID)
	if err != nil {
		return nil, "", err
	}

	apikey, rawToken, err := pri.CreateAPIKey(name)
	if err != nil {
		return nil, "", err
	}

	err = cmd.SaveAuthAPIKey(ctx, apikey, pri)
	if err != nil {
		log.Println("save error on SaveAuthAPIKey:", err)
		return nil, "", err
	}

	tx.Commit()

	return apikey, rawToken, err
}

func (u *Administration) ListRoles(ctx Context) ([]*membership.Role, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, systemError("Could not start transaction", err)
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	manager := membership.NewManager(q)

	ids, err := manager.EnumerateRoleIDs(ctx)
	if err != nil {
		return nil, err
	}

	list, err := manager.FindRoles(ctx, ids)
	if err != nil {
		return nil, err
	}

	return list, nil
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
