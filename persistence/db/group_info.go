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

type GroupRow struct {
	SeqID       int64  `db:"seq_id,primary,auto_increment"`
	GroupID     string `db:"group_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func groupRowFromEntity(g *entity.Group) *GroupRow {
	return &GroupRow{
		GroupID:     g.GroupID.String(),
		Name:        g.Name,
		Description: g.Description,
	}
}

func (s *Service) findGroupsByPrincipalSeqID(ctx Context, conn seacle.Selectable, priSeqIDs []int64) (map[int64][]*entity.Group, error) {
	priGroupMaps := []*PrincipalGroupMapRow{}
	err := seacle.Select(ctx, conn, &priGroupMaps, `WHERE principal_seq_id IN (?)`, priSeqIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(priGroupMaps) == 0 {
		return map[int64][]*entity.Group{}, nil
	}

	groupSeqIDMap := map[int64]int64{} // RoleSeqID -> PrincipalSeqID map
	groupSeqIDs := []int64{}
	for _, v := range priGroupMaps {
		groupSeqIDMap[v.GroupSeqID] = v.PrincipalSeqID
		groupSeqIDs = append(groupSeqIDs, v.GroupSeqID)
	}

	groupMap, err := s.findGroups(ctx, conn, groupSeqIDs)
	if err != nil {
		return nil, err
	}

	result := map[int64][]*entity.Group{}
	for groupSeqID, r := range groupMap {
		principalSeqID := groupSeqIDMap[groupSeqID]
		result[principalSeqID] = append(result[principalSeqID], r)
	}

	return result, nil
}

func (s *Service) SaveGroup(ctx context.Context, conn seacle.Executable, g *entity.Group) error {
	row := &GroupRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE group_id = ?", g.GroupID.String())
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}
	if err == sql.ErrNoRows {
		err = s.createGroupRow(ctx, conn, g)
		if err != nil {
			log.Println(err)
			return err
		}
		err = seacle.SelectRow(ctx, conn, row, "WHERE group_id = ?", g.GroupID.String())
	} else {
		err = s.updateGroupRow(ctx, conn, g, row)
	}

	if err != nil {
		log.Println(err)
		return err
	}

	err = s.saveGroupRole(ctx, conn, g, row)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Service) DeleteGroup(ctx context.Context, conn seacle.Executable, g *entity.Group) error {
	row := &GroupRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE group_id = ?", g.GroupID.String())
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}
	if err == sql.ErrNoRows {
		return nil
	} else {
		err = s.deleteGroupRow(ctx, conn, g, row)
	}

	// TODO delete group-role

	return err
}

func (s *Service) createGroupRow(ctx context.Context, conn seacle.Executable, g *entity.Group) error {
	row := groupRowFromEntity(g)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updateGroupRow(ctx context.Context, conn seacle.Executable, g *entity.Group, row *GroupRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock role row: err=%s", err)
	}

	// update row
	row.Name = g.Name
	row.Description = g.Description
	err = seacle.Update(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update role row: err=%s", err)
	}

	return nil
}

func (s *Service) deleteGroupRow(ctx context.Context, conn seacle.Executable, g *entity.Group, row *GroupRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock role row: err=%s", err)
	}

	// delete row
	err = seacle.Delete(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to delete role row: err=%s", err)
	}

	return nil
}

func (s *Service) saveGroupRole(ctx Context, conn seacle.Executable, g *entity.Group, gRow *GroupRow) error {
	roles := g.Roles()
	roleSeqIDs := make([]int64, 0, len(roles))
	if len(roles) != 0 {
		roleIDs := make([]string, 0, len(roles))
		for _, v := range roles {
			roleIDs = append(roleIDs, v.RoleID.String())
		}
		roleRows := []*RoleRow{}
		err := seacle.Select(ctx, conn, &roleRows, `WHERE role_id IN (?)`, roleIDs)
		if err != nil {
			return err
		}
		for _, v := range roleRows {
			roleSeqIDs = append(roleSeqIDs, v.SeqID)
		}
	}

	groupRoleMaps := []*GroupRoleMapRow{}
	err := seacle.Select(ctx, conn, &groupRoleMaps, `WHERE group_seq_id = ?`, gRow.SeqID)
	if err != nil {
		return err
	}
	args := make([]relationRow, 0, len(groupRoleMaps))
	for _, v := range groupRoleMaps {
		args = append(args, v)
	}

	added, deleted := compareSeqID(roleSeqIDs, args)
	if len(added) != 0 {
		for _, roleSeqID := range added {
			_, err = seacle.Insert(ctx, conn, &GroupRoleMapRow{
				GroupSeqID: gRow.SeqID,
				RoleSeqID:  roleSeqID,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	if len(deleted) != 0 {
		for _, roleSeqID := range deleted {
			err = seacle.Delete(ctx, conn, &GroupRoleMapRow{
				GroupSeqID: gRow.SeqID,
				RoleSeqID:  roleSeqID,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (s *Service) FindGroups(ctx Context, conn seacle.Selectable, groupIDs []string) ([]*entity.Group, error) {
	if len(groupIDs) == 0 {
		return []*entity.Group{}, nil
	}

	groupRows := make([]*GroupRow, 0, 8)
	err := seacle.Select(ctx, conn, &groupRows, `WHERE group_id IN (?) ORDER BY seq_id`, groupIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// seq_idを抽出
	groupSeqIDs := make([]int64, 0, len(groupRows))
	for _, v := range groupRows {
		groupSeqIDs = append(groupSeqIDs, v.SeqID)
	}
	if len(groupSeqIDs) == 0 {
		return []*entity.Group{}, nil
	}

	groupMap, err := s.findGroups(ctx, conn, groupSeqIDs)
	if err != nil {
		return nil, err
	}
	groups := []*entity.Group{}
	for _, v := range groupSeqIDs {
		if g, exist := groupMap[v]; exist {
			groups = append(groups, g)
		}
	}
	return groups, nil
}

func (s *Service) EnumerateGroupIDs(ctx Context, conn seacle.Selectable) ([]string, error) {
	groups := make([]*GroupRow, 0, 8)
	err := seacle.Select(ctx, conn, &groups, "ORDER BY seq_id")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]string, 0, len(groups))
	for _, v := range groups {
		result = append(result, v.GroupID)
	}

	return result, nil
}

// findRoles returns SeqID -> *entity.Group map
func (s *Service) findGroups(ctx Context, conn seacle.Selectable, groupSeqIDs []int64) (map[int64]*entity.Group, error) {
	if len(groupSeqIDs) == 0 {
		return map[int64]*entity.Group{}, nil
	}

	// roles
	groupRoleMaps := []*GroupRoleMapRow{}
	roleMap := map[int64][]*entity.Role{} // GroupSeqID -> []Role map
	err := seacle.Select(ctx, conn, &groupRoleMaps, `WHERE group_seq_id IN (?)`, groupSeqIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(groupRoleMaps) != 0 {
		roleSeqIDMap := map[int64]int64{} // RoleSeqID -> GroupSeqID map
		roleSeqIDs := []int64{}
		for _, v := range groupRoleMaps {
			roleSeqIDMap[v.RoleSeqID] = v.GroupSeqID
			roleSeqIDs = append(roleSeqIDs, v.RoleSeqID)
		}

		roles, err := s.findRoles(ctx, conn, roleSeqIDs)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for roleSeqID, r := range roles {
			groupSeqID := roleSeqIDMap[roleSeqID]
			roleMap[groupSeqID] = append(roleMap[groupSeqID], r)
		}
	}

	groupRows := []*GroupRow{}
	groupMap := map[int64]*entity.Group{}
	err = seacle.Select(ctx, conn, &groupRows, `WHERE seq_id IN (?) ORDER BY seq_id`, groupSeqIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	f := entity.NewFactory(s.q)
	for _, v := range groupRows {
		r := f.NewGroup(uuid.MustParse(v.GroupID), v.Name, v.Description, roleMap[v.SeqID])
		groupMap[v.SeqID] = r
	}

	return groupMap, nil
}
