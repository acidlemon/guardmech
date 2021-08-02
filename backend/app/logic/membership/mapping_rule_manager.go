package membership

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"sort"
	"strings"

	"github.com/acidlemon/guardmech/backend/app/logic/auth"
)

type MappingRuleManager struct {
	rules    []*MappingRule
	inquirer *auth.GroupInquirer
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

func NewMappingRuleManager(rules []*MappingRule, inquirer *auth.GroupInquirer) *MappingRuleManager {
	s := make(MappingRuleSlice, 0, len(rules))
	s = append(s, rules...)
	sort.Stable(s)

	return &MappingRuleManager{
		rules:    s,
		inquirer: inquirer,
	}
}

func (m *MappingRuleManager) FindMatchedRules(ctx context.Context, token *auth.OpenIDToken) ([]*MappingRule, error) {
	addr, err := mail.ParseAddress(token.Email)
	if err != nil {
		return nil, fmt.Errorf("malformed email address: address=%s, err=%s", token.Email, err)
	}

	result := make([]*MappingRule, 0, len(m.rules))

	for _, r := range m.rules {
		switch r.RuleType {
		case MappingSpecificDomain:
			if strings.HasSuffix(addr.Address, "@"+r.Detail) {
				result = append(result, r)
			}

		case MappingWholeDomain:
			if strings.HasSuffix(addr.Address, r.Detail) {
				result = append(result, r)
			}

		case MappingGroupMember:
			if m.inquirer != nil {
				ok, err := m.inquirer.IsMember(ctx, addr.Address, r.Detail)
				if err != nil {
					log.Println("failed to inquire member:", err)
					continue
				}
				if ok {
					result = append(result, r)
				}
			}

		case MappingSpecificAddress:
			if addr.Address == r.Detail {
				result = append(result, r)
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Matched rule was not found")
	}

	return result, nil
}
