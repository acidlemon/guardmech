package membership

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Group struct {
	SeqID       int64
	UniqueID    uuid.UUID `json:"unique_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (r *Group) AttachPermission(ctx context.Context, conn *sql.Conn, pe *Permission) error {
	_, err := conn.ExecContext(ctx, `INSERT INTO group_permission_map (group_id, permission_id) VALUES (?, ?)`, r.SeqID, pe.SeqID)
	return err
}
