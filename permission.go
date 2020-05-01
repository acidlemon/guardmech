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

func CreatePermission(ctx context.Context, tx *sql.Tx, name string) (*Permission, error) {
	res, err := tx.ExecContext(ctx, `INSERT INTO permission (name, description) VALUES (?, ?)`, name, "")
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

func scanPermissionRow(r RowScanner) (*Permission, error) {
	var id int64
	var name, description string
	err := r.Scan(&id, &name, &description)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &Permission{
		ID:          id,
		Name:        name,
		Description: description,
	}, nil
}
