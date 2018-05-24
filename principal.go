package guardmech

import (
	"context"
	"database/sql"
	"log"

	"github.com/k0kubun/pp"
)

type Auth struct {
	ID        int64
	Account   string
	Principal *Principal
}

type APIKey struct {
	ID        int64
	Token     string
	Principal *Principal
}

type Principal struct {
	ID          int64
	Name        string
	Description string
}

func FindPrincipal(ctx context.Context, conn *sql.Conn, account string) (*Principal, error) {
	row := conn.QueryRowContext(ctx, `SELECT p.id, p.name, p.description FROM auth AS a JOIN principal AS p ON a.principal_id = p.id where a.account = ?`, account)
	var id int64
	var name, description string
	err := row.Scan(&id, &name, &description)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &Principal{
		ID:          id,
		Name:        name,
		Description: description,
	}, nil
}

func HasPrincipal(ctx context.Context, conn *sql.Conn) (bool, error) {
	row := conn.QueryRowContext(ctx, `SELECT COUNT(*) AS cnt FROM principal`)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return false, err
	}

	if cnt == 0 {
		return false, nil
	}
	return true, nil
}

func CreateFirstPrincipal(ctx context.Context, conn *sql.Conn, account string) error {
	tx, err := conn.BeginTx(ctx, nil)
	commited := false
	defer func() {
		if !commited {
			tx.Rollback()
		}
	}()

	pr, err := CreateAuthPrincipal(ctx, conn, account)
	if err != nil {
		tx.Rollback()
		return err
	}

	pe, err := CreatePermission(ctx, conn, PermissionOwner)
	if err != nil {
		tx.Rollback()
		return err
	}

	r, err := CreateRole(ctx, conn, RoleOwner)
	if err != nil {
		tx.Rollback()
		return err
	}

	r.AttachPermission(ctx, conn, pe)
	pr.AttachRole(ctx, conn, r)

	tx.Commit()

	return nil
}

func CreateAuthPrincipal(ctx context.Context, conn *sql.Conn, account string) (*Principal, error) {
	// insert to principal
	res, err := conn.ExecContext(ctx, `INSERT INTO principal (name, description) VALUES (?, ?)`, account, "")
	if err != nil {
		return nil, err
	}
	principalID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// insert to auth
	_, err = conn.ExecContext(ctx, `INSERT INTO auth (account, principal_id) VALUES (?, ?)`, account, principalID)
	if err != nil {
		return nil, err
	}

	return &Principal{
		ID:          principalID,
		Name:        account,
		Description: "",
	}, nil
}

func (pr *Principal) AttachRole(ctx context.Context, conn *sql.Conn, r *Role) error {
	_, err := conn.ExecContext(ctx, `INSERT INTO principal_role_map (principal_id, role_id) VALUES (?, ?)`, pr.ID, r.ID)
	return err
}

func (pr *Principal) FindRole(ctx context.Context, conn *sql.Conn) ([]*Role, error) {
	result := make([]*Role, 0, 32)

	pp.Print(pr)

	// find role (direct attached)
	rows, err := conn.QueryContext(ctx, `SELECT r.id, r.name, r.description FROM principal_role_map AS m JOIN role_info AS r ON m.role_id = r.id WHERE m.principal_id = ?`, pr.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			log.Println("scan error:", err)
			return nil, err
		}
		result = append(result, &Role{
			ID:          id,
			Name:        name,
			Description: description,
		})
	}

	// find role (attached via group)
	// TODO

	return result, nil
}

//func (pr *Principal)
