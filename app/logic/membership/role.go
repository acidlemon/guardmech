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
	err := r.AttachPermission(perm)
	if err != nil {
		return nil, err
	}
	return perm, nil
}

func (r *Role) AttachPermission(p *Permission) error {
	for _, v := range r.permissions {
		if v.PermissionID == p.PermissionID {
			// already exists
			return nil // do nothing, no error (idempotence)
		}
	}

	r.permissions = append(r.permissions, p)
	return nil
}

func (r *Role) DetachPermission(p *Permission) error {
	for i, v := range r.permissions {
		if v.PermissionID == p.PermissionID {
			r.permissions = append(r.permissions[:i], r.permissions[i+1:]...)
			return nil
		}
	}
	return nil // Not Found, but no error (idempotence)
}
