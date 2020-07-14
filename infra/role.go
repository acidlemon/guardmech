package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
	"github.com/acidlemon/seacle"
)

type Role membership.Role

func (m *Membership) FetchAllRole(ctx context.Context, conn *sql.Conn) ([]*membership.Role, error) {
	roles := []*Role{}
	err := seacle.Select(ctx, conn, &roles, "")
	if err != nil {
		return nil, err
	}

	result := make([]*membership.Role, 0, len(roles))
	for _, v := range roles {
		r := membership.Role(*v)
		result = append(result, &r)
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
	// new Role
	role := Role(*r)
	seqID, err := seacle.Insert(ctx, tx, &role)
	if err != nil {
		return 0, err
	}

	return seqID, nil
}

func (s *Membership) updateRole(ctx context.Context, tx *db.Tx, r *membership.Role) (int64, error) {
	roleForUpdate := &Role{}
	err := seacle.SelectRow(ctx, tx, roleForUpdate, `WHERE unique_id = ?`, r.UniqueID)
	if err != nil {
		// TODO: fallback to createRole
		return 0, nil
	}

	if roleForUpdate.SeqID != r.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, r.SeqID=%d", roleForUpdate.SeqID, r.SeqID)
	}

	// lock row
	err = seacle.SelectRow(ctx, tx, roleForUpdate, `FOR UPDATE WHERE seq_id = ?`, roleForUpdate.SeqID)
	if err != nil {
		return 0, fmt.Errorf("failed to lock role row: err=%s", err)
	}

	// update row
	role := Role(*r)
	err = seacle.Update(ctx, tx, &role)
	if err != nil {
		return 0, fmt.Errorf("failed to update role row: err=%s", err)
	}

	return r.SeqID, nil
}

func uniqRoles(roles []*Role) []*Role {
	m := map[int64]struct{}{}
	result := make([]*Role, 0, len(roles))

	for _, r := range roles {
		if _, exist := m[r.SeqID]; !exist {
			result = append(result, r)
			m[r.SeqID] = struct{}{}
		}
	}

	return result
}

func roleIDs(roles []*Role) []int64 {
	result := make([]int64, 0, len(roles))
	for _, r := range roles {
		result = append(result, r.SeqID)
	}
	return result
}
