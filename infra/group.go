package infra

import (
	"context"
	"fmt"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
	"github.com/acidlemon/seacle"
)

type Group membership.Group

func (s *Membership) SaveGroup(ctx context.Context, tx *db.Tx, gr *membership.Group) (int64, error) {
	if gr.SeqID == 0 {
		return s.createGroup(ctx, tx, gr)
	} else {
		return s.updateGroup(ctx, tx, gr)
	}
}

func (s *Membership) createGroup(ctx context.Context, tx *db.Tx, gr *membership.Group) (int64, error) {
	// new Group
	group := Group(*gr)
	seqID, err := seacle.Insert(ctx, tx, &group)
	if err != nil {
		return 0, err
	}

	return seqID, nil
}

func (s *Membership) updateGroup(ctx context.Context, tx *db.Tx, gr *membership.Group) (int64, error) {
	groupForUpdate := &Group{}
	err := seacle.SelectRow(ctx, tx, groupForUpdate, `WHERE unique_id = ?`, gr.UniqueID)
	if err != nil {
		// TODO: fallback to createRole
		return 0, nil
	}

	if groupForUpdate.SeqID != gr.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, gr.SeqID=%d", groupForUpdate.SeqID, gr.SeqID)
	}

	// lock row
	err = seacle.SelectRow(ctx, tx, groupForUpdate, `FOR UPDATE WHERE seq_id = ?`, groupForUpdate.SeqID)
	if err != nil {
		return 0, fmt.Errorf("failed to lock role row: err=%s", err)
	}

	// update row
	group := Group(*gr)
	err = seacle.Update(ctx, tx, &group)
	if err != nil {
		return 0, fmt.Errorf("failed to update role row: err=%s", err)
	}

	return gr.SeqID, nil
}
