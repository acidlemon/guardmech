package membership

import (
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

/*
func (r *Group) AttachPermission(ctx Context, conn *sql.Conn, pe *Permission) error {
	_, err := conn.ExecContext(ctx, `INSERT INTO group_permission_map (group_id, permission_id) VALUES (?, ?)`, r.SeqID, pe.SeqID)
	return err
}
*/
