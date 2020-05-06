package membership

import (
	"context"

	"github.com/acidlemon/guardmech/db"
	"github.com/google/uuid"
)

const (
	RoleOwner = "Guardmech-Owner"
)

type Role struct {
	SeqID       int64
	UniqueID    uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (r *Role) AttachPermission(ctx context.Context, tx *db.Tx, pe *Permission) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO role_permission_map (role_id, permission_id) VALUES (?, ?)`, r.SeqID, pe.SeqID)
	return err
}
