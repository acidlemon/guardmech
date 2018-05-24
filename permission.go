package guardmech

import (
	"context"
	"database/sql"
)

const (
	PermissionOwner = "_GUARDMECH_OWNER"
)

type Permission struct {
	ID          int64
	Name        string
	Description string
}

func CreatePermission(ctx context.Context, conn *sql.Conn, name string) (*Permission, error) {
	res, err := conn.ExecContext(ctx, `INSERT INTO permission (name, description) VALUES (?, ?)`, name, "")
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Permission{
		ID:   id,
		Name: name,
	}, nil
}
