package db

import (
	"log"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type MappingRuleRow struct {
	SeqID         int64  `db:"seq_id,primary,auto_increment"`
	MappingRuleID string `db:"mapping_rule_id"`
	Type          int    `db:"type"`
	Detail        string `db:"detail"`
	Name          string `db:"name"`
	Description   string `db:"description"`
}

func mappingRuleRowFromEntity(m *entity.MappingRule) *MappingRuleRow {
	return &MappingRuleRow{
		MappingRuleID: m.MappingRuleID.String(),
		Type:          int(m.Type),
		Detail:        m.Detail,
		Name:          m.Name,
		Description:   m.Description,
	}
}

func (m *MappingRuleRow) ToEntity() *entity.MappingRule {
	return &entity.MappingRule{
		MappingRuleID: uuid.MustParse(m.MappingRuleID),
		Type:          entity.MappingType(m.Type),
		Detail:        m.Detail,
		Name:          m.Name,
		Description:   m.Description,
	}
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

	mrs := []*entity.MappingRule{}
	for _, v := range mrRows {
		mrs = append(mrs, v.ToEntity())
	}
	return mrs, nil
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
