package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
)

func (m *Membership) FetchAllRole(ctx context.Context, conn *sql.Conn) ([]*membership.Role, error) {
	rows, err := conn.QueryContext(ctx, `SELECT r.id, r.name, r.description FROM role_info AS r`)
	if err != nil {
		return nil, err
	}

	result := []*membership.Role{}
	for rows.Next() {
		r, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

func (s *Membership) SaveRole(ctx context.Context, tx *db.Tx, r *membership.Role) (int64, error) {
	if r.SeqID == 0 {
		return s.createRole(ctx, tx, r)
	} else {
		return s.updateRole(ctx, tx, r)
	}
}

func (s *Membership) createRole(ctx context.Context, tx *db.Tx, r *membership.Role) (int64, error) {
	// new Principal
	res, err := tx.ExecContext(ctx,
		`INSERT INTO permission (unique_id, name, description) VALUES (?, ?, ?)`,
		r.UniqueID, r.Name, r.Description,
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

func (s *Membership) updateRole(ctx context.Context, tx *db.Tx, r *membership.Role) (int64, error) {
	var seqID int64
	row := tx.QueryRowContext(ctx,
		`SELECT seq_id FROM role_info WHERE unique_id = ?`,
		r.UniqueID,
	)
	err := row.Scan(&seqID)
	if err != nil {
		// TODO: fallback to createRole?
		return 0, nil
	}

	if seqID != r.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, r.SeqID=%d", seqID, r.SeqID)
	}

	// lock row
	row = tx.QueryRowContext(ctx,
		`SELECT seq_id FROM role_info FOR UPDATE WHERE seq_id = ?`, seqID)
	err = row.Scan(&seqID)
	if err != nil {
		return 0, nil
	}

	// update row
	//tx.ExecContext(ctx, `UPDATE`)

	return seqID, nil
}

func uniqRoles(roles []*membership.Role) []*membership.Role {
	m := map[int64]struct{}{}
	result := make([]*membership.Role, 0, len(roles))

	for _, r := range roles {
		if _, exist := m[r.SeqID]; !exist {
			result = append(result, r)
			m[r.SeqID] = struct{}{}
		}
	}

	return result
}

func roleIDs(roles []*membership.Role) []int64 {
	result := make([]int64, 0, len(roles))
	for _, r := range roles {
		result = append(result, r.SeqID)
	}
	return result
}
