package infra

import (
	"context"
	"fmt"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
)

type Permission membership.Permission

func (s *Membership) SavePermission(ctx context.Context, tx *db.Tx, pe *membership.Permission) (int64, error) {
	if pe.SeqID == 0 {
		return s.createPermission(ctx, tx, pe)
	} else {
		return s.updatePermission(ctx, tx, pe)
	}
}

func (s *Membership) createPermission(ctx context.Context, tx *db.Tx, pe *membership.Permission) (int64, error) {
	// new Principal
	res, err := tx.ExecContext(ctx,
		`INSERT INTO permission (unique_id, name, description) VALUES (?, ?, ?)`,
		pe.UniqueID, pe.Name, pe.Description,
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

func (s *Membership) updatePermission(ctx context.Context, tx *db.Tx, pe *membership.Permission) (int64, error) {
	var seqID int64
	row := tx.QueryRowContext(ctx,
		`SELECT seq_id FROM permission WHERE unique_id = ?`,
		pe.UniqueID,
	)
	err := row.Scan(&seqID)
	if err != nil {
		// TODO: fallback to createPermission?
		return 0, nil
	}

	if seqID != pe.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, pe.SeqID=%d", seqID, pe.SeqID)
	}

	// lock row
	row = tx.QueryRowContext(ctx,
		`SELECT seq_id FROM permission FOR UPDATE WHERE seq_id = ?`, seqID)
	err = row.Scan(&seqID)
	if err != nil {
		return 0, nil
	}

	// update row
	//tx.ExecContext(ctx, `UPDATE`)

	return seqID, nil
}
