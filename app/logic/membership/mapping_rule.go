package membership

import "github.com/google/uuid"

type MappingType int

type MappingRule struct {
	MappingRuleID uuid.UUID
	Type          MappingType
	Detail        string
	Name          string
	Description   string
}
