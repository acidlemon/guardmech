package usecase

import (
	"fmt"
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

func (u *Administration) DeletePrincipal(ctx Context, id string) error {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	pris, err := manager.FindPrincipals(ctx, []string{id})
	if err != nil {
		return err
	}

	if len(pris) == 0 {
		return fmt.Errorf("principal not found")
	}

	pri := pris[0]
	err = cmd.DeletePrincipal(ctx, pri)
	if err != nil {
		log.Println("failed to delete principal")
		return err
	}

	tx.Commit()

	return nil
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
func (u *Administration) DeleteRole(ctx Context, id string) error {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	roles, err := manager.FindRoles(ctx, []string{id})
	if err != nil {
		return err
	}

	if len(roles) == 0 {
		return fmt.Errorf("role not found")
	}

	r := roles[0]
	err = cmd.DeleteRole(ctx, r)
	if err != nil {
		log.Println("failed to delete role")
		return err
	}

	tx.Commit()

	return nil
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
func (u *Administration) FetchGroup(ctx Context, id string) (*membership.Group, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	q := persistence.NewQuery(conn)
	manager := membership.NewManager(q)

	groups, err := manager.FindGroups(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("group not found")
	}

	return groups[0], nil
}
func (u *Administration) UpdateGroup(ctx Context) (*membership.Group, error) {
	return nil, nil
}
func (u *Administration) DeleteGroup(ctx Context, id string) error {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	groups, err := manager.FindGroups(ctx, []string{id})
	if err != nil {
		return err
	}

	if len(groups) == 0 {
		return fmt.Errorf("group not found")
	}

	g := groups[0]
	err = cmd.DeleteGroup(ctx, g)
	if err != nil {
		log.Println("failed to delete group")
		return err
	}

	tx.Commit()

	return nil
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
func (u *Administration) DeletePermission(ctx Context, id string) error {
	conn, tx, err := db.GetTxConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer tx.AutoRollback()

	q := persistence.NewQuery(tx)
	cmd := persistence.NewCommand(tx)
	manager := membership.NewManager(q)

	perms, err := manager.FindPermissions(ctx, []string{id})
	if err != nil {
		return err
	}

	if len(perms) == 0 {
		return fmt.Errorf("permission not found")
	}

	perm := perms[0]
	err = cmd.DeletePermission(ctx, perm)
	if err != nil {
		log.Println("failed to delete permission")
		return err
	}

	tx.Commit()

	return nil
}
