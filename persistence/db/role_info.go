package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type RoleRow struct {
	SeqID       int64  `db:"seq_id,primary,auto_increment"`
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

// func (r *RoleRow) ToEntity() *entity.Role {
// 	return &entity.Role{
// 		RoleID:      uuid.MustParse(r.RoleID),
// 		Name:        r.Name,
// 		Description: r.Description,
// 	}
// }

func (s *Service) findRolesByPrincipalSeqID(ctx Context, conn seacle.Selectable, priSeqIDs []int64) (map[int64][]*entity.Role, error) {
	priRoleMaps := []*PrincipalRoleMapRow{}
	err := seacle.Select(ctx, conn, &priRoleMaps, `WHERE principal_seq_id IN (?)`, priSeqIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(priRoleMaps) == 0 {
		return map[int64][]*entity.Role{}, nil
	}

	roleSeqIDMap := map[int64]int64{} // RoleSeqID -> PrincipalSeqID map
	roleSeqIDs := []int64{}
	for _, v := range priRoleMaps {
		roleSeqIDMap[v.RoleSeqID] = v.PrincipalSeqID
		roleSeqIDs = append(roleSeqIDs, v.RoleSeqID)
	}

	roleMap, err := s.findRoles(ctx, conn, roleSeqIDs)
	if err != nil {
		return nil, err
	}

	result := map[int64][]*entity.Role{}
	for roleSeqID, r := range roleMap {
		principalSeqID := roleSeqIDMap[roleSeqID]
		result[principalSeqID] = append(result[principalSeqID], r)
	}

	return result, nil
}

func (s *Service) SaveRole(ctx context.Context, conn seacle.Executable, r *entity.Role) error {
	row := &RoleRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE role_id = ?", r.RoleID.String())
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	if err == sql.ErrNoRows {
		err = s.createRoleRow(ctx, conn, r)
		if err != nil {
			log.Println(err)
			return err
		}
		err = seacle.SelectRow(ctx, conn, row, "WHERE role_id = ?", r.RoleID.String())
	} else {
		err = s.updateRoleRow(ctx, conn, r, row)
	}
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.saveRolePermission(ctx, conn, r, row)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Service) DeleteRole(ctx context.Context, conn seacle.Executable, r *entity.Role) error {
	row := &RoleRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE role_id = ?", r.RoleID.String())
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	if err == sql.ErrNoRows {
		return nil
	} else {
		err = s.deleteRoleRow(ctx, conn, r, row)
	}

	// TODO delete RolePermission

	return err
}

func (s *Service) createRoleRow(ctx context.Context, conn seacle.Executable, r *entity.Role) error {
	row := roleRowFromEntity(r)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updateRoleRow(ctx context.Context, conn seacle.Executable, r *entity.Role, row *RoleRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
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

func (s *Service) deleteRoleRow(ctx context.Context, conn seacle.Executable, r *entity.Role, row *RoleRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock role row: err=%s", err)
	}

	// delete row
	err = seacle.Delete(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update role row: err=%s", err)
	}

	return nil
}

func (s *Service) saveRolePermission(ctx Context, conn seacle.Executable, r *entity.Role, roleRow *RoleRow) error {
	perms := r.Permissions()
	permSeqIDs := make([]int64, 0, len(perms))
	if len(perms) != 0 {
		permIDs := make([]string, 0, len(perms))
		for _, v := range perms {
			permIDs = append(permIDs, v.PermissionID.String())
		}
		permRows := []*PermissionRow{}
		err := seacle.Select(ctx, conn, &permRows, `WHERE permission_id in (?)`, permIDs)
		if err != nil {
			return err
		}
		for _, v := range permRows {
			permSeqIDs = append(permSeqIDs, v.SeqID)
		}
	}

	rolePermMaps := []*RolePermissionMapRow{}
	err := seacle.Select(ctx, conn, &rolePermMaps, `WHERE role_seq_id = ?`, roleRow.SeqID)
	if err != nil {
		return err
	}
	args := make([]relationRow, 0, len(rolePermMaps))
	for _, v := range rolePermMaps {
		args = append(args, v)
	}

	added, deleted := compareSeqID(permSeqIDs, args)
	if len(added) != 0 {
		for _, permSeqID := range added {
			_, err = seacle.Insert(ctx, conn, &RolePermissionMapRow{
				RoleSeqID:       roleRow.SeqID,
				PermissionSeqID: permSeqID,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	if len(deleted) != 0 {
		for _, permSeqID := range deleted {
			err = seacle.Delete(ctx, conn, &RolePermissionMapRow{
				RoleSeqID:       roleRow.SeqID,
				PermissionSeqID: permSeqID,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}

func (s *Service) FindRoles(ctx Context, conn seacle.Selectable, roleIDs []string) ([]*entity.Role, error) {
	if len(roleIDs) == 0 {
		return []*entity.Role{}, nil
	}

	roleRows := make([]*RoleRow, 0, 8)
	err := seacle.Select(ctx, conn, &roleRows, `WHERE role_id IN (?)`, roleIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// seq_idを抽出
	roleSeqIDs := make([]int64, 0, len(roleRows))
	for _, v := range roleRows {
		roleSeqIDs = append(roleSeqIDs, v.SeqID)
	}
	if len(roleSeqIDs) == 0 {
		return []*entity.Role{}, nil
	}

	roleMap, err := s.findRoles(ctx, conn, roleSeqIDs)
	if err != nil {
		return nil, err
	}
	roles := []*entity.Role{}
	for _, v := range roleMap {
		roles = append(roles, v)
	}
	return roles, nil
}

func (s *Service) EnumerateRoleIDs(ctx Context, conn seacle.Selectable) ([]string, error) {
	roles := make([]*RoleRow, 0, 8)
	err := seacle.Select(ctx, conn, &roles, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]string, 0, len(roles))
	for _, v := range roles {
		result = append(result, v.RoleID)
	}

	return result, nil
}

// findRoles returns SeqID -> *entity.Role map
func (s *Service) findRoles(ctx Context, conn seacle.Selectable, roleSeqIDs []int64) (map[int64]*entity.Role, error) {
	if len(roleSeqIDs) == 0 {
		return map[int64]*entity.Role{}, nil
	}

	// permissions
	rolePermMaps := []*RolePermissionMapRow{}
	permMap := map[int64][]*entity.Permission{} // RoleSeqID -> []Permission map
	err := seacle.Select(ctx, conn, &rolePermMaps, `WHERE role_seq_id IN (?)`, roleSeqIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(rolePermMaps) != 0 {
		permSeqIDMap := map[int64]int64{} // PermissionSeqID -> RoleSeqID map
		permSeqIDs := []int64{}
		for _, v := range rolePermMaps {
			permSeqIDMap[v.PermissionSeqID] = v.RoleSeqID
			permSeqIDs = append(permSeqIDs, v.PermissionSeqID)
		}

		perms := []*PermissionRow{}
		err = seacle.Select(ctx, conn, &perms, `WHERE seq_id IN (?)`, permSeqIDs)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for _, v := range perms {
			roleSeqID := permSeqIDMap[v.SeqID]
			permMap[roleSeqID] = append(permMap[roleSeqID], v.ToEntity())
		}
	}

	roleRows := []*RoleRow{}
	roleMap := map[int64]*entity.Role{}
	err = seacle.Select(ctx, conn, &roleRows, `WHERE seq_id IN (?)`, roleSeqIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	f := entity.NewFactory(s.q)
	for _, v := range roleRows {
		r := f.NewRole(uuid.MustParse(v.RoleID), v.Name, v.Description, permMap[v.SeqID])
		roleMap[v.SeqID] = r
	}

	return roleMap, nil
}
