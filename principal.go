package guardmech

import (
	"context"
	"database/sql"
	"log"

	"github.com/k0kubun/pp"
)

type Auth struct {
	ID        int64      `json:"id"`
	Account   string     `json:"account"`
	Principal *Principal `json:"-"`
}

type APIKey struct {
	ID        int64      `json:"id"`
	Token     string     `json:"token"`
	Principal *Principal `json:"-"`
}

type Principal struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PrincipalPayload struct {
	Principal *Principal `json:"principal"`
	Auths     []*Auth    `json:"auths"`
	APIKeys   []*APIKey  `json:"api_keys"`
}

func FindPrincipal(ctx context.Context, conn *sql.Conn, account string) (*Principal, error) {
	row := conn.QueryRowContext(ctx, `SELECT p.id, p.name, p.description FROM auth AS a JOIN principal AS p ON a.principal_id = p.id where a.account = ?`, account)
	pr, err := scanPrincipalRow(row)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func FindPrincipalByID(ctx context.Context, conn *sql.Conn, id string) (*PrincipalPayload, error) {
	row := conn.QueryRowContext(ctx, `SELECT p.id, p.name, p.description FROM principal p WHERE id = ?`, id)
	pr, err := scanPrincipalRow(row)
	if err != nil {
		return nil, err
	}

	auths := make([]*Auth, 0, 4)
	rows, err := conn.QueryContext(ctx, `SELECT a.id, a.account FROM auth AS a WHERE a.principal_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a, err := scanAuthRow(rows, pr)
		if err != nil {
			return nil, err
		}
		auths = append(auths, a)
	}

	apikeys := make([]*APIKey, 0, 4)
	rows, err = conn.QueryContext(ctx, `SELECT a.id, a.Token FROM api_key AS a WHERE a.principal_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a, err := scanAPIKeyRow(rows, pr)
		if err != nil {
			return nil, err
		}
		apikeys = append(apikeys, a)
	}

	return &PrincipalPayload{
		Principal: pr,
		Auths:     auths,
		APIKeys:   apikeys,
	}, nil
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

func scanPrincipalRow(r RowScanner) (*Principal, error) {
	var id int64
	var name, description string
	err := r.Scan(&id, &name, &description)
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

func scanAuthRow(r RowScanner, pr *Principal) (*Auth, error) {
	var id int64
	var account string
	err := r.Scan(&id, &account)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &Auth{
		ID:        id,
		Account:   account,
		Principal: pr,
	}, nil
}

func scanAPIKeyRow(r RowScanner, pr *Principal) (*APIKey, error) {
	var id int64
	var token string
	err := r.Scan(&id, &token)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &APIKey{
		ID:        id,
		Token:     token,
		Principal: pr,
	}, nil
}

func FetchAllPrincipal(ctx context.Context, conn *sql.Conn) ([]*Principal, error) {
	rows, err := conn.QueryContext(ctx, `SELECT p.id, p.name, p.description FROM principal p`)
	if err != nil {
		return nil, err
	}

	result := []*Principal{}
	for rows.Next() {
		pr, err := scanPrincipalRow(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, pr)
	}

	return result, nil
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
