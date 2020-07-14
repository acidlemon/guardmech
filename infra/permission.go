package infra

import (
	"context"
	"fmt"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
	"github.com/acidlemon/seacle"
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
	// new Permission
	perm := Permission(*pe)
	seqID, err := seacle.Insert(ctx, tx, &perm)
	if err != nil {
		return 0, err
	}

	return seqID, nil
}

func (s *Membership) updatePermission(ctx context.Context, tx *db.Tx, pe *membership.Permission) (int64, error) {
	permForUpdate := &Permission{}
	err := seacle.SelectRow(ctx, tx, permForUpdate, `WHERE unique_id = ?`, pe.UniqueID)
	if err != nil {
		// TODO: fallback to createRole
		return 0, nil
	}

	if permForUpdate.SeqID != pe.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, r.SeqID=%d", permForUpdate.SeqID, pe.SeqID)
	}

	// lock row
	err = seacle.SelectRow(ctx, tx, permForUpdate, `FOR UPDATE WHERE seq_id = ?`, permForUpdate.SeqID)
	if err != nil {
		return 0, fmt.Errorf("failed to lock permission row: err=%s", err)
	}

	// update row
	perm := Permission(*pe)
	err = seacle.Update(ctx, tx, &perm)
	if err != nil {
		return 0, fmt.Errorf("failed to update permission row: err=%s", err)
	}

	return pe.SeqID, nil
}
