package guardmech

import (
	"context"
	"database/sql"
)

type Group struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateGroup(ctx context.Context, conn *sql.Conn, name string) (*Group, error) {
	res, err := conn.ExecContext(ctx, `INSERT INTO group_info (name, description) VALUES (?, ?)`, name, "")
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Group{
		ID:   id,
		Name: name,
	}, nil
}

func (r *Group) AttachPermission(ctx context.Context, conn *sql.Conn, pe *Permission) error {
	_, err := conn.ExecContext(ctx, `INSERT INTO group_permission_map (group_id, permission_id) VALUES (?, ?)`, r.ID, pe.ID)
	return err
}

func scanGroupRow(r RowScanner) (*Group, error) {
	var id int64
	var name, description string
	err := r.Scan(&id, &name, &description)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &Group{
		ID:          id,
		Name:        name,
		Description: description,
	}, nil
}
