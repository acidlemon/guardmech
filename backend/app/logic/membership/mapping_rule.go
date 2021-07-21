package membership

import "github.com/google/uuid"

type MappingType int

//specific domain, 2=whole domain, 3=member of 4=specific address

const (
	MappingSpecificDomain MappingType = iota + 1
	MappingWholeDomain
	MappingGroupMember
	MappingSpecificAddress
)

type MappingRule struct {
	MappingRuleID uuid.UUID
	RuleType      MappingType
	Detail        string
	Name          string
	Description   string
	Priority      int

	associatedGroup *Group
	associatedRole  *Role
}

func (m *MappingRule) AssociatedGroup() *Group {
	return m.associatedGroup
}

func (m *MappingRule) AssociatedRole() *Role {
	return m.associatedRole
}
