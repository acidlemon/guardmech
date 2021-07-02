package db

import (
	entity "github.com/acidlemon/guardmech/app/logic/membership"
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
