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

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	pri, err := manager.CreatePrincipal(ctx, name, description)
	if err != nil {
		return nil, err
	}

	err = cmd.SavePrincipal(ctx, pri)
	if err != nil {
		log.Println("failed to save new principal")
		return nil, err
	}

	tx.Commit()

	return pri, nil
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

func (u *Administration) CreateRole(ctx Context, name, description string) (*membership.Role, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	log.Println("name=", name, "description=", description)

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	r, err := manager.CreateRole(ctx, name, description)
	if err != nil {
		return nil, err
	}

	err = cmd.SaveRole(ctx, r)
	if err != nil {
		log.Println("failed to save new role")
		return nil, err
	}

	tx.Commit()

	return r, nil
}
func (u *Administration) FetchRole(ctx Context, id string) (*membership.Role, error) {
	return nil, nil
}
func (u *Administration) UpdateRole(ctx Context, id string) (*membership.Role, error) {
	return nil, nil
}
func (u *Administration) ListMappingRules(ctx Context) ([]*membership.MappingRule, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, systemError("Could not start transaction", err)
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	manager := membership.NewManager(q)

	ids, err := manager.EnumerateMappingRuleIDs(ctx)
	if err != nil {
		return nil, err
	}

	list, err := manager.FindMappingRules(ctx, ids)
	if err != nil {
		return nil, err
	}

	return list, nil
}
func (u *Administration) CreateMappingRule(ctx Context) (*membership.MappingRule, error) {
	return nil, nil
}
func (u *Administration) FetchMappingRule(ctx Context, id string) (*membership.MappingRule, error) {
	return nil, nil
}
func (u *Administration) UpdateMappingRule(ctx Context, id string) (*membership.MappingRule, error) {
	return nil, nil
}
func (u *Administration) ListGroups(ctx Context) ([]*membership.Group, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, systemError("Could not start transaction", err)
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	manager := membership.NewManager(q)

	ids, err := manager.EnumerateGroupIDs(ctx)
	if err != nil {
		return nil, err
	}

	list, err := manager.FindGroups(ctx, ids)
	if err != nil {
		return nil, err
	}

	return list, nil
}
func (u *Administration) CreateGroup(ctx Context, name, description string) (*membership.Group, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	g, err := manager.CreateGroup(ctx, name, description)
	if err != nil {
		return nil, err
	}

	err = cmd.SaveGroup(ctx, g)
	if err != nil {
		log.Println("failed to save new group")
		return nil, err
	}

	tx.Commit()

	return g, nil
}
func (u *Administration) FetchGroup(ctx Context) (*membership.Group, error) {
	return nil, nil
}
func (u *Administration) UpdateGroup(ctx Context) (*membership.Group, error) {
	return nil, nil
}
func (u *Administration) ListPermissions(ctx Context) ([]*membership.Permission, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, systemError("Could not start transaction", err)
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	manager := membership.NewManager(q)

	ids, err := manager.EnumeratePermissionIDs(ctx)
	if err != nil {
		return nil, err
	}

	list, err := manager.FindPermissions(ctx, ids)
	if err != nil {
		return nil, err
	}

	return list, nil
}
func (u *Administration) CreatePermission(ctx Context, name, description string) (*membership.Permission, error) {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	perm, err := manager.CreatePermission(ctx, name, description)
	if err != nil {
		return nil, err
	}

	err = cmd.SavePermission(ctx, perm)
	if err != nil {
		log.Println("failed to save new group")
		return nil, err
	}

	tx.Commit()

	return perm, nil
}
func (u *Administration) FetchPermission(ctx Context) (*membership.Permission, error) {
	return nil, nil
}
func (u *Administration) UpdatePermission(ctx Context) (*membership.Permission, error) {
	return nil, nil
}
