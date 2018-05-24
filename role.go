package guardmech

import (
	"context"
	"database/sql"
)

const (
	RoleOwner = "Guardmech-Owner"
)

type Role struct {
	ID          int64
	Name        string
	Description string
}

func CreateRole(ctx context.Context, conn *sql.Conn, name string) (*Role, error) {
	res, err := conn.ExecContext(ctx, `INSERT INTO role_info (name, description) VALUES (?, ?)`, name, "")
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Role{
		ID:   id,
		Name: name,
	}, nil
}

func (r *Role) AttachPermission(ctx context.Context, conn *sql.Conn, pe *Permission) error {
	_, err := conn.ExecContext(ctx, `INSERT INTO role_permission_map (role_id, permission_id) VALUES (?, ?)`, r.ID, pe.ID)
	return err
}
