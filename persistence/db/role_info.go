package db

import (
	"context"
	"database/sql"
	"fmt"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/seacle"
)

type RoleRow struct {
	SeqID       int64  `db:"seq_id,primary"`
	RoleID      string `db:"role_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func roleRowFromEntity(r *entity.Role) *RoleRow {
	return &RoleRow{
		RoleID:      r.RoleID.String(),
		Name:        r.Name,
		Description: r.Description,
	}
}

/*
func (m *Service) FetchAllRole(ctx context.Context, conn *sql.Conn) ([]*membership.Role, error) {
	roles := []*RoleRow{}
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
*/

func (s *Service) SaveRole(ctx context.Context, conn seacle.Executable, r *entity.Role) error {
	row := &RoleRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE role_id = ?", r.RoleID.String())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return s.createRole(ctx, conn, r)
	} else {
		return s.updateRole(ctx, conn, r, row)
	}
}

func (s *Service) createRole(ctx context.Context, conn seacle.Executable, r *entity.Role) error {
	row := roleRowFromEntity(r)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updateRole(ctx context.Context, conn seacle.Executable, r *entity.Role, row *RoleRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `FOR UPDATE WHERE seq_id = ?`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock role row: err=%s", err)
	}

	// update row
	row.Name = r.Name
	row.Description = r.Description
	err = seacle.Update(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update role row: err=%s", err)
	}

	return nil
}

/*
func uniqRoles(roles []*RoleRow) []*RoleRow {
	m := map[int64]struct{}{}
	result := make([]*RoleRow, 0, len(roles))

	for _, r := range roles {
		if _, exist := m[r.SeqID]; !exist {
			result = append(result, r)
			m[r.SeqID] = struct{}{}
		}
	}

	return result
}

func roleIDs(roles []*RoleRow) []int64 {
	result := make([]int64, 0, len(roles))
	for _, r := range roles {
		result = append(result, r.SeqID)
	}
	return result
}
*/
