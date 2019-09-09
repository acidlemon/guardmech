package guardmech

import (
	"context"
	"database/sql"
)

const (
	RoleOwner = "Guardmech-Owner"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

func scanRoleRow(r RowScanner) (*Role, error) {
	var id int64
	var name, description string
	err := r.Scan(&id, &name, &description)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &Role{
		ID:          id,
		Name:        name,
		Description: description,
	}, nil
}

func FetchAllRole(ctx context.Context, conn *sql.Conn) ([]*Role, error) {
	rows, err := conn.QueryContext(ctx, `SELECT r.id, r.name, r.description FROM role AS r`)
	if err != nil {
		return nil, err
	}

	result := []*Role{}
	for rows.Next() {
		r, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}
