package infra

import (
	"context"
	"fmt"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
)

type Context = context.Context

type Auth membership.Auth

func (s *Membership) SaveAuth(ctx Context, tx *db.Tx, a *membership.Auth) (int64, error) {
	if a.SeqID == 0 {
		return s.createAuth(ctx, tx, a)
	} else {
		return s.updateAuth(ctx, tx, a)
	}
}

func (s *Membership) createAuth(ctx Context, tx *db.Tx, a *membership.Auth) (int64, error) {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO auth (unique_id, issuer, subject, email, principal_id)
			VALUES (?, ?, ?, ?, ?)`,
		a.UniqueID, a.Issuer, a.Subject, a.Email, a.Principal.SeqID,
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

func (s *Membership) updateAuth(ctx Context, tx *db.Tx, a *membership.Auth) (int64, error) {
	var seqID int64
	row := tx.QueryRowContext(ctx,
		`SELECT seq_id FROM auth WHERE unique_id = ?`,
		a.UniqueID,
	)
	err := row.Scan(&seqID)
	if err != nil {
		// TODO: fallback to createAuth?
		return 0, nil
	}

	if seqID != a.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, a.SeqID=%d", seqID, a.SeqID)
	}

	// lock row
	row = tx.QueryRowContext(ctx,
		`SELECT seq_id FROM auth FOR UPDATE WHERE seq_id = ?`, seqID)
	err = row.Scan(&seqID)
	if err != nil {
		return 0, nil
	}

	// update row
	//tx.ExecContext(ctx, `UPDATE`)

	return seqID, nil
}
