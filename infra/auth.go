package infra

import (
	"context"
	"fmt"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
	"github.com/acidlemon/seacle"
)

type Context = context.Context

type Auth struct {
	membership.Auth
	PrincipalID int64
}

func (s *Membership) SaveAuth(ctx Context, tx *db.Tx, a *membership.Auth) (int64, error) {
	if a.SeqID == 0 {
		return s.createAuth(ctx, tx, a)
	} else {
		return s.updateAuth(ctx, tx, a)
	}
}

func (s *Membership) createAuth(ctx Context, tx *db.Tx, a *membership.Auth) (int64, error) {
	au := Auth{
		*a, a.Principal.SeqID,
	}
	seqID, err := seacle.Insert(ctx, tx, &au)
	if err != nil {
		return 0, err
	}
	return seqID, nil
}

func (s *Membership) updateAuth(ctx Context, tx *db.Tx, a *membership.Auth) (int64, error) {
	auForUpdate := &Auth{}
	err := seacle.SelectRow(ctx, tx, auForUpdate, `WHERE unique_id = ?`, a.UniqueID)
	if err != nil {
		// TODO: fallback to createAuth?
		return 0, nil
	}

	if auForUpdate.SeqID != a.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, a.SeqID=%d", auForUpdate.SeqID, a.SeqID)
	}

	// lock row
	err = seacle.SelectRow(ctx, tx, auForUpdate, `FOR UPDATE WHERE seq_id = ?`, auForUpdate.SeqID)
	if err != nil {
		return 0, fmt.Errorf("failed to lock auth row: err=%s", err)
	}

	// update row
	au := Auth{
		*a, a.Principal.SeqID,
	}
	err = seacle.Update(ctx, tx, &au)
	if err != nil {
		return 0, fmt.Errorf("failed to update auth row: err=%s", err)
	}

	return a.SeqID, nil
}
