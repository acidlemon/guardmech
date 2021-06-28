package membership

import (
	"github.com/google/uuid"
)

type Group struct {
	GroupID     uuid.UUID
	Name        string
	Description string

	roles []*Role
}

func (g *Group) Roles() []*Role {
	if g.roles == nil {
		return []*Role{}
	}
	return g.roles
}

/*
func (r *Group) AttachPermission(ctx Context, conn *sql.Conn, pe *Permission) error {
	_, err := conn.ExecContext(ctx, `INSERT INTO group_permission_map (group_id, permission_id) VALUES (?, ?)`, r.SeqID, pe.SeqID)
	return err
}
*/
