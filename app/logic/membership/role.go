package membership

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	RoleOwnerName        = "_GuardmechOwnerRole"
	RoleOwnerDescription = "Owner principal of guardmech"
	RoleOwnerID          = "b8cc3e1a-867e-4c2d-b163-c9feb5683388"
)

type Role struct {
	RoleID      uuid.UUID
	Name        string
	Description string

	permissions []*Permission
}

func newRole(name, description string) *Role {
	roleID := uuid.New()
	if name == RoleOwnerName {
		roleID = uuid.MustParse(RoleOwnerID)
	}

	r := &Role{
		RoleID:      roleID,
		Name:        name,
		Description: description,
	}

	return r
}

func (r *Role) Permissions() []*Permission {
	if r.permissions == nil {
		return []*Permission{}
	}

	return r.permissions
}

func (r *Role) AttachNewPermission(ctx Context, name, description string) (*Permission, error) {
	if name == "" {
		return nil, fmt.Errorf("AttachNewPermission: name is required")
	}

	perm := newPermission(name, description)
	r.permissions = append(r.permissions, perm)
	return perm, nil
}

/*
func (r *Role) AttachPermission(ctx Context, tx *db.Tx, pe *Permission) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO role_permission_map (role_id, permission_id) VALUES (?, ?)`, r.SeqNo, pe.SeqNo)
	return err
}
*/
