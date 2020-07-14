package infra

import (
	"fmt"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
	"github.com/acidlemon/seacle"
)

type APIKey struct {
	membership.APIKey
	PrincipalID int64
}

func (s *Membership) SaveAPIKey(ctx Context, tx *db.Tx, a *membership.APIKey) (int64, error) {
	if a.SeqID == 0 {
		return s.createAPIKey(ctx, tx, a)
	} else {
		return s.updateAPIKey(ctx, tx, a)
	}
}

func (s *Membership) createAPIKey(ctx Context, tx *db.Tx, a *membership.APIKey) (int64, error) {
	ap := APIKey{
		*a, a.Principal.SeqID,
	}
	seqID, err := seacle.Insert(ctx, tx, &ap)
	if err != nil {
		return 0, err
	}
	return seqID, nil
}

func (s *Membership) updateAPIKey(ctx Context, tx *db.Tx, a *membership.APIKey) (int64, error) {
	apForUpdate := &APIKey{}
	err := seacle.SelectRow(ctx, tx, apForUpdate, `WHERE unique_id = ?`, a.UniqueID)
	if err != nil {
		// TODO: fallback to createAPIKey?
		return 0, nil
	}

	if apForUpdate.SeqID != a.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, a.SeqID=%d", apForUpdate.SeqID, a.SeqID)
	}

	// lock row
	err = seacle.SelectRow(ctx, tx, apForUpdate, `FOR UPDATE WHERE seq_id = ?`, apForUpdate.SeqID)
	if err != nil {
		return 0, fmt.Errorf("failed to lock api_key row: err=%s", err)
	}

	// update row
	ap := APIKey{
		*a, a.Principal.SeqID,
	}
	err = seacle.Update(ctx, tx, &ap)
	if err != nil {
		return 0, fmt.Errorf("failed to update api_key row: err=%s", err)
	}

	return a.SeqID, nil
}
