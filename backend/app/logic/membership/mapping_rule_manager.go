package membership

import (
	"fmt"
	"net/mail"
	"sort"
	"strings"

	"github.com/acidlemon/guardmech/app/logic/auth"
)

type MappingRuleManager struct {
	rules []*MappingRule
}

// for sort.Interface
type MappingRuleSlice []*MappingRule

func (s MappingRuleSlice) Len() int {
	return len(s)
}
func (s MappingRuleSlice) Less(i, j int) bool {
	if s[i].Priority < s[j].Priority {
		return true
	}
	return false
}
func (s MappingRuleSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NewMappingRuleManager(rules []*MappingRule) *MappingRuleManager {
	s := make(MappingRuleSlice, 0, len(rules))
	s = append(s, rules...)
	sort.Stable(s)

	return &MappingRuleManager{
		rules: s,
	}
}

func (m *MappingRuleManager) FindMatchedRule(token *auth.OpenIDToken) (*MappingRule, error) {
	addr, err := mail.ParseAddress(token.Email)
	if err != nil {
		return nil, fmt.Errorf("malformed email address: address=%s, err=%s", token.Email, err)
	}

	for _, r := range m.rules {
		switch r.RuleType {
		case MappingSpecificDomain:
			if strings.HasSuffix(addr.Address, "@"+r.Detail) {
				return r, nil
			}

		case MappingWholeDomain:
			if strings.HasSuffix(addr.Address, r.Detail) {
				return r, nil
			}

		case MappingGroupMember:
			// TODO つくる
			continue

		case MappingSpecificAddress:
			if addr.Address == r.Detail {
				return r, nil
			}
		}
	}

	return nil, fmt.Errorf("Matched rule was not found")
}

//func (m *MappingRuleManager)
