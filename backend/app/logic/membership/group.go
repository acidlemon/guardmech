package membership

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	GroupOwnerName        = "_GuardmechOwnerGroup"
	GroupOwnerDescription = "Owner group of guardmech"
	GroupOwnerID          = "6f43787e-1a18-42dc-86dc-78c81c681bda"
)

type Group struct {
	GroupID     uuid.UUID
	Name        string
	Description string

	roles []*Role
}

func newGroup(name, description string) *Group {
	groupID := uuid.New()
	if name == GroupOwnerName {
		groupID = uuid.MustParse(PermissionOwnerID)
	}

	return &Group{
		GroupID:     groupID,
		Name:        name,
		Description: description,
	}
}

func (g *Group) Roles() []*Role {
	if g.roles == nil {
		return []*Role{}
	}
	return g.roles
}

func (g *Group) AttachNewRole(name, description string) (*Role, error) {
	if name == "" {
		return nil, fmt.Errorf("AttachNewRole: name is required")
	}

	// create New Role
	r := newRole(name, description)
	err := g.AttachRole(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (g *Group) AttachRole(r *Role) error {
	for _, v := range g.roles {
		if v.RoleID == r.RoleID {
			// already exists
			return nil // TODO error?
		}
	}

	g.roles = append(g.roles, r)

	return nil
}

func (g *Group) DetachRole(r *Role) error {
	for i, v := range g.roles {
		if v.RoleID == r.RoleID {
			g.roles = append(g.roles[:i], g.roles[i+1:]...)
			return nil
		}
	}
	return nil // Not Found, but no error (idempotence)
}

func (g *Group) HavingPermissions() []*Permission {
	roles := g.Roles()

	tmp := []*Permission{}
	if len(roles) == 0 {
		return tmp
	}

	for _, r := range roles {
		tmp = append(tmp, r.Permissions()...)
	}

	// uniq list
	existsMap := map[string]bool{}
	result := make([]*Permission, 0, len(tmp))
	for _, p := range tmp {
		uidstr := p.PermissionID.String()
		if _, exist := existsMap[uidstr]; !exist {
			result = append(result, p)
			existsMap[uidstr] = true
		}
	}

	return result
}
