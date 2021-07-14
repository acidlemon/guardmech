package db

import (
	"database/sql"
	"fmt"
	"log"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

const AssociatedTypeGroup = 1
const AssociatedTypeRole = 2

type MappingRuleRow struct {
	SeqID           int64  `db:"seq_id,primary,auto_increment"`
	MappingRuleID   string `db:"mapping_rule_id"`
	RuleType        int    `db:"rule_type"`
	Detail          string `db:"detail"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	Priority        int    `db:"priority"`
	AssociationType int    `db:"association_type"`
	AssociationID   string `db:"association_id"`
}

func mappingRuleRowFromEntity(m *entity.MappingRule) *MappingRuleRow {
	row := &MappingRuleRow{
		MappingRuleID: m.MappingRuleID.String(),
		RuleType:      int(m.RuleType),
		Detail:        m.Detail,
		Name:          m.Name,
		Description:   m.Description,
		Priority:      m.Priority,
	}

	g := m.AssociatedGroup()
	if g != nil {
		row.AssociationType = AssociatedTypeGroup
		row.AssociationID = g.GroupID.String()
	}

	r := m.AssociatedRole()
	if r != nil {
		row.AssociationType = AssociatedTypeRole
		row.AssociationID = r.RoleID.String()
	}

	return row
}

func (s *Service) FindMappingRules(ctx Context, conn seacle.Selectable, mappingRuleIDs []string) ([]*entity.MappingRule, error) {
	if len(mappingRuleIDs) == 0 {
		return []*entity.MappingRule{}, nil
	}

	mrRows := make([]*MappingRuleRow, 0, 8)
	err := seacle.Select(ctx, conn, &mrRows, `WHERE mapping_rule_id IN (?)`, mappingRuleIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(mrRows) == 0 {
		return []*entity.MappingRule{}, nil
	}

	roleIDs := []string{}
	groupIDs := []string{}
	for _, v := range mrRows {
		if v.AssociationType == AssociatedTypeGroup {
			groupIDs = append(groupIDs, v.AssociationID)
		} else if v.AssociationType == AssociatedTypeRole {
			roleIDs = append(roleIDs, v.AssociationID)
		}
	}

	// Associated Group
	groups, err := s.FindGroups(ctx, conn, groupIDs)
	if err != nil {
		return nil, err
	}

	// Associated Role
	roles, err := s.FindRoles(ctx, conn, roleIDs)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.MappingRule, 0, len(mrRows))
	f := entity.NewFactory(s.q)
	for _, v := range mrRows {
		var group *entity.Group
		var role *entity.Role
		if v.AssociationType == AssociatedTypeGroup {
			for _, w := range groups {
				if w.GroupID.String() == v.AssociationID {
					group = w
					break
				}
			}
		} else if v.AssociationType == AssociatedTypeRole {
			for _, w := range roles {
				if w.RoleID.String() == v.AssociationID {
					role = w
					break
				}
			}
		}

		result = append(result, f.NewMappingRule(
			uuid.MustParse(v.MappingRuleID), entity.MappingType(v.RuleType), v.Detail,
			v.Name, v.Description, v.Priority, group, role,
		))
	}
	return result, nil
}

func (s *Service) EnumerateMappingRuleIDs(ctx Context, conn seacle.Selectable) ([]string, error) {
	mrs := make([]*MappingRuleRow, 0, 8)
	err := seacle.Select(ctx, conn, &mrs, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]string, 0, len(mrs))
	for _, v := range mrs {
		result = append(result, v.MappingRuleID)
	}

	return result, nil
}

func (s *Service) SaveMappingRule(ctx Context, conn seacle.Executable, rule *entity.MappingRule) error {
	row := &MappingRuleRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE mapping_rule_id = ?", rule.MappingRuleID.String())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return s.createMappingRuleRow(ctx, conn, rule)
	} else {
		return s.updateMappingRuleRow(ctx, conn, rule, row)
	}
}

func (s *Service) DeleteMappingRule(ctx Context, conn seacle.Executable, rule *entity.MappingRule) error {
	row := &MappingRuleRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE mapping_rule_id = ?", rule.MappingRuleID.String())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return nil // do nothing
	}

	return s.deleteMappingRuleRow(ctx, conn, rule, row)
}

func (s *Service) createMappingRuleRow(ctx Context, conn seacle.Executable, rule *entity.MappingRule) error {
	row := mappingRuleRowFromEntity(rule)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updateMappingRuleRow(ctx Context, conn seacle.Executable, rule *entity.MappingRule, row *MappingRuleRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock mapping_rule row: err=%s", err)
	}

	// update row
	row.Name = rule.Name
	row.Description = rule.Description
	err = seacle.Update(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update mapping_rule row: err=%s", err)
	}

	return nil
}

func (s *Service) deleteMappingRuleRow(ctx Context, conn seacle.Executable, rule *entity.MappingRule, row *MappingRuleRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock mapping_rule row: err=%s", err)
	}

	// delete row
	err = seacle.Delete(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to delete mapping_rule row: err=%s", err)
	}

	return nil
}
