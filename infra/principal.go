package infra

//go:generate go run github.com/acidlemon/seacle/cmd/seacle

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
	"github.com/acidlemon/seacle"
	"github.com/jmoiron/sqlx"
)

//+table: principal
type Principal membership.Principal

type Membership struct {
}

func (s *Membership) HasPrincipal(ctx context.Context, conn *sql.Conn) (bool, error) {
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

func (s *Membership) FetchPrincipalPayload(ctx context.Context, conn *sql.Conn, id int64) (*membership.PrincipalPayload, error) {

	// row := conn.QueryRowContext(ctx,
	// 	`SELECT p.seq_id, p.unique_id, p.name, p.description FROM principal AS p WHERE seq_id = ?`, id)
	// pr, err := scanPrincipalRow(row)
	// if err != nil {
	// 	return nil, err
	// }

	pr := &Principal{}
	err := seacle.SelectRow(ctx, conn, pr, `WHERE seq_id = ?`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	mempr := membership.Principal(*pr)

	auths := make([]*membership.Auth, 0, 4)
	rows, err := conn.QueryContext(ctx,
		`SELECT a.seq_id, a.unique_id, a.issuer, a.subject, a.email FROM auth AS a WHERE a.principal_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a, err := scanAuthRow(rows)
		if err != nil {
			return nil, err
		}
		a.Principal = &mempr
		auths = append(auths, a)
	}

	apikeys := make([]*membership.APIKey, 0, 4)
	rows, err = conn.QueryContext(ctx,
		`SELECT a.seq_id, a.unique_id, a.token FROM api_key AS a WHERE a.principal_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a, err := scanAPIKeyRow(rows)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		a.Principal = &mempr
		apikeys = append(apikeys, a)
	}

	groups := make([]*membership.Group, 0, 8)
	groupIDs := make([]int64, 0, 8)
	rows, err = conn.QueryContext(ctx,
		`SELECT r.seq_id, r.unique_id, r.name, r.description FROM principal_group_map AS m
		JOIN group_info AS r ON m.group_id = r.seq_id WHERE m.principal_id = ?`, id)
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
		grp := membership.Group(*r)
		groups = append(groups, &grp)
		groupIDs = append(groupIDs, r.SeqID)
	}

	roles := make([]*membership.Role, 0, 8)
	rows, err = conn.QueryContext(ctx,
		`SELECT r.seq_id, r.unique_id, r.name, r.description FROM principal_role_map AS m 
		JOIN role_info AS r ON m.role_id = r.seq_id WHERE m.principal_id = ?`, id)
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
		query, args, err := sqlx.In(`SELECT r.seq_id, r.unique_id, r.name, r.description FROM group_role_map AS m 
		JOIN role_info AS r ON m.role_id = r.seq_id WHERE m.group_id IN (?))`, groupIDs)
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

	permissions := make([]*membership.Permission, 0, 8)

	if len(roles) > 0 {
		query, args, err := sqlx.In(`SELECT pe.seq_id, pe.unique_id, pe.name, pe.description FROM role_permission_map AS m 
		JOIN permission AS pe ON m.permission_id = pe.seq_id WHERE m.role_id IN (?)`, roleIDs(roles))
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
	}

	return &membership.PrincipalPayload{
		Principal:   &mempr,
		Auths:       auths,
		APIKeys:     apikeys,
		Groups:      groups,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func (s *Membership) SavePrincipal(ctx context.Context, tx *db.Tx, pri *membership.Principal) (int64, error) {
	if pri.SeqID == 0 {
		return s.createPrincipal(ctx, tx, pri)
	} else {
		return s.updatePrincipal(ctx, tx, pri)
	}
}

func (s *Membership) createPrincipal(ctx context.Context, tx *db.Tx, pri *membership.Principal) (int64, error) {
	// new Principal
	res, err := tx.ExecContext(ctx,
		`INSERT INTO principal (unique_id, name, description) VALUES (?, ?, ?)`,
		pri.UniqueID, pri.Name, pri.Description,
	)
	if err != nil {
		return 0, err
	}
	seqID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return seqID, nil
}

func (s *Membership) updatePrincipal(ctx context.Context, tx *db.Tx, pri *membership.Principal) (int64, error) {
	var seqID int64
	row := tx.QueryRowContext(ctx,
		`SELECT seq_id FROM principal WHERE unique_id = ?`,
		pri.UniqueID,
	)
	err := row.Scan(&seqID)
	if err != nil {
		// TODO: fallback to createPrincipal?
		return 0, nil
	}

	if seqID != pri.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, pri.SeqID=%d", seqID, pri.SeqID)
	}

	// lock row
	row = tx.QueryRowContext(ctx,
		`SELECT p.seq_id FROM principal AS p FOR UPDATE WHERE p.seq_id = ?`, seqID)
	err = row.Scan(&seqID)
	if err != nil {
		return 0, nil
	}

	// update row
	//tx.ExecContext(ctx, `UPDATE`)

	return seqID, nil
}

func (s *Membership) FindPrincipal(ctx context.Context, conn *sql.Conn, issuer, subject string) (*membership.Principal, error) {
	// 	row := conn.QueryRowContext(ctx, `SELECT p.seq_id, p.unique_id, p.name, p.description FROM auth AS a
	//   JOIN principal AS p ON a.principal_id = p.seq_id
	//   WHERE a.issuer = ? AND a.subject = ?`, issuer, subject)
	// 	pr, err := scanPrincipalRow(row)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return pr, nil

	pr := &Principal{}
	err := seacle.SelectRow(ctx, conn, pr,
		`JOIN auth AS a ON a.principal_id = principal.seq_id WHERE a.issuer = ? AND a.subject = ?`,
		issuer, subject)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	mempr := membership.Principal(*pr)

	return &mempr, nil
}

func (s *Membership) FetchAllPrincipal(ctx context.Context, conn *sql.Conn) ([]*membership.Principal, error) {
	rows, err := conn.QueryContext(ctx, `SELECT p.seq_id, p.unique_id, p.name, p.description FROM principal p`)
	if err != nil {
		return nil, err
	}

	result := []*membership.Principal{}
	for rows.Next() {
		pr, err := scanPrincipalRow(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, pr)
	}

	return result, nil
}
