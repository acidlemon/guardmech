package guardmech

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/k0kubun/pp"
)

type Auth struct {
	ID        int64      `json:"id"`
	Issuer    string     `json:"issuer"`
	Subject   string     `json:"subject"`
	Email     string     `json:"email"`
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

type PrincipalDetailed struct {
	Principal
	Groups []string
	Roles  []string
}

type PrincipalPayload struct {
	Principal   *Principal    `json:"principal"`
	Auths       []*Auth       `json:"auths"`
	APIKeys     []*APIKey     `json:"api_keys"`
	Groups      []*Group      `json:"groups"`
	Roles       []*Role       `json:"roles"`
	Permissions []*Permission `json:"permissions"`
}

func FindPrincipal(ctx context.Context, conn *sql.Conn, issuer, subject string) (*Principal, error) {
	log.Println(issuer)
	log.Println(subject)
	row := conn.QueryRowContext(ctx, `SELECT p.id, p.name, p.description FROM auth AS a 
  JOIN principal AS p ON a.principal_id = p.id
  WHERE a.issuer = ? AND a.subject = ?`, issuer, subject)
	pr, err := scanPrincipalRow(row)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func FetchPrincipalPayload(ctx context.Context, conn *sql.Conn, id int64) (*PrincipalPayload, error) {
	row := conn.QueryRowContext(ctx,
		`SELECT p.id, p.name, p.description FROM principal AS p WHERE id = ?`, id)
	pr, err := scanPrincipalRow(row)
	if err != nil {
		return nil, err
	}

	auths := make([]*Auth, 0, 4)
	rows, err := conn.QueryContext(ctx,
		`SELECT a.id, a.issuer, a.subject, a.email FROM auth AS a WHERE a.principal_id = ?`, id)
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
	rows, err = conn.QueryContext(ctx,
		`SELECT a.id, a.Token FROM api_key AS a WHERE a.principal_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a, err := scanAPIKeyRow(rows, pr)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		apikeys = append(apikeys, a)
	}

	groups := make([]*Group, 0, 8)
	groupIDs := make([]int64, 0, 8)
	rows, err = conn.QueryContext(ctx,
		`SELECT r.id, r.name, r.description FROM principal_group_map AS m
		JOIN group_info AS r ON m.group_id = r.id WHERE m.principal_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r, err := scanGroupRow(rows)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		groups = append(groups, r)
		groupIDs = append(groupIDs, r.ID)
	}

	roles := make([]*Role, 0, 8)
	rows, err = conn.QueryContext(ctx,
		`SELECT r.id, r.name, r.description FROM principal_role_map AS m 
		JOIN role_info AS r ON m.role_id = r.id WHERE m.principal_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r, err := scanRoleRow(rows)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		roles = append(roles, r)
	}

	if len(groupIDs) > 0 {
		query, args, err := sqlx.In(`SELECT r.id, r.name, r.description FROM group_role_map AS m 
		JOIN role_info AS r ON m.role_id = r.id WHERE m.group_id IN (?))`, groupIDs)
		rows, err = conn.QueryContext(ctx, query, args...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			r, err := scanRoleRow(rows)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			roles = append(roles, r)
		}

		roles = uniqRoles(roles)
	}

	permissions := make([]*Permission, 0, 8)
	query, args, err := sqlx.In(`SELECT p.id, p.name, p.description FROM role_permission_map AS m 
	JOIN permission AS p ON m.permission_id = p.id WHERE m.role_id IN (?)`, roleIDs(roles))
	rows, err = conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p, err := scanPermissionRow(rows)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return &PrincipalPayload{
		Principal:   pr,
		Auths:       auths,
		APIKeys:     apikeys,
		Groups:      groups,
		Roles:       roles,
		Permissions: permissions,
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
	var issuer, subject, email string
	err := r.Scan(&id, &issuer, &subject, &email)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &Auth{
		ID:        id,
		Issuer:    issuer,
		Subject:   subject,
		Email:     email,
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

	//result := []*PrincipalDetailed{}
	result := []*Principal{}
	for rows.Next() {
		pr, err := scanPrincipalRow(rows)
		if err != nil {
			return nil, err
		}

		//result = append(result, &PrincipalDetailed{Principal: pr})
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

func CreateFirstPrincipal(ctx context.Context, conn *sql.Conn, idToken *OpenIDToken) error {
	tx, err := conn.BeginTx(ctx, nil)
	commited := false
	defer func() {
		if !commited {
			tx.Rollback()
		}
	}()

	pr, err := CreateAuthPrincipal(ctx, conn, idToken)
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

func CreateAuthPrincipal(ctx context.Context, conn *sql.Conn, idToken *OpenIDToken) (*Principal, error) {
	// insert to principal
	name := idToken.Name
	if name == "" {
		name = idToken.Email
	}
	res, err := conn.ExecContext(ctx, `INSERT INTO principal (name, description) VALUES (?, ?)`, name, "")
	if err != nil {
		return nil, err
	}
	principalID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// insert to auth
	_, err = conn.ExecContext(ctx, `INSERT INTO auth (issuer, subject, email, principal_id) VALUES (?, ?, ?, ?)`,
		idToken.Issuer, idToken.Sub, idToken.Email, principalID)
	if err != nil {
		return nil, err
	}

	return &Principal{
		ID:          principalID,
		Name:        name,
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
