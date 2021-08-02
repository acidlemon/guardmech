package membership

import (
	"fmt"
	"log"
	"strings"

	"github.com/acidlemon/guardmech/backend/app/logic/auth"
	"github.com/google/uuid"
)

type Manager struct {
	q Query
}

func NewManager(q Query) *Manager {
	return &Manager{
		q: q,
	}
}

func (s *Manager) validateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	if strings.Contains(name, ";") {
		return fmt.Errorf("name cannot contains ';'")
	}

	return nil
}

func (s *Manager) FindPrincipalByID(ctx Context, principalID string) (*Principal, error) {
	pri, err := s.q.FindPrincipals(ctx, []string{principalID})
	if err != nil {
		return nil, err
	}
	if len(pri) == 0 {
		log.Println("FindPrincipalByID:", principalID, "no entry")
		return nil, ErrNoEntry
	}

	return pri[0], nil
}

func (s *Manager) FindPrincipalByOIDC(ctx Context, issuer, subject string) (*Principal, error) {
	pri, err := s.q.FindPrincipalIDByOIDC(ctx, issuer, subject)
	if err != nil {
		return nil, err
	}

	return pri, nil
}

func (s *Manager) EnumeratePrincipalIDs(ctx Context) ([]string, error) {
	IDs, err := s.q.EnumeratePrincipalIDs(ctx)
	if err != nil {
		return nil, err
	}

	if len(IDs) == 0 {
		return []string{}, nil
	}

	result := make([]string, 0, len(IDs))
	for _, u := range IDs {
		result = append(result, u.String())
	}

	return result, nil
}

func (m *Manager) FindPrincipals(ctx Context, ids []string) ([]*Principal, error) {
	list, err := m.q.FindPrincipals(ctx, ids)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Manager) CreatePrincipalFromOpenID(ctx Context, token *auth.OpenIDToken) (*Principal, *OIDCAuthorization, error) {
	a := &OIDCAuthorization{
		OIDCAuthID: uuid.New(),
		Issuer:     token.Issuer,
		Subject:    token.Sub,
		Email:      token.Email,
		Name:       token.Name,
	}
	pri := &Principal{
		PrincipalID: uuid.New(),
		Name:        token.Name,
		Description: token.Email,
		auth:        a,
	}

	return pri, a, nil
}

func (s *Manager) CreatePrincipalFromAPIKey(ctx Context, name, apiKey string) (*Principal, error) {
	if err := s.validateName(name); err != nil {
		return nil, fmt.Errorf("CreatePrincipalFromAPIKey: %s", err.Error())
	}

	//	s.cmd.CreatePrincipal()
	return &Principal{}, nil
}

func (s *Manager) CreatePrincipal(ctx Context, name, description string) (*Principal, error) {
	if err := s.validateName(name); err != nil {
		return nil, fmt.Errorf("CreatePrincipal: %s", err.Error())
	}

	return &Principal{
		PrincipalID: uuid.New(),
		Name:        name,
		Description: description,
	}, nil
}

func (s *Manager) CreatePrincipalFromRules(ctx Context, token *auth.OpenIDToken, rules []*MappingRule) (*Principal, *OIDCAuthorization, error) {
	pri, a, err := s.CreatePrincipalFromOpenID(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	for _, rule := range rules {
		g := rule.AssociatedGroup()
		if g != nil {
			pri.AttachGroup(g)
		}

		r := rule.AssociatedRole()
		if r != nil {
			pri.AttachRole(r)
		}
	}

	return pri, a, nil
}

func (s *Manager) CreateRole(ctx Context, name, description string) (*Role, error) {
	if err := s.validateName(name); err != nil {
		return nil, fmt.Errorf("CreateRole: %s", err.Error())
	}

	return &Role{
		RoleID:      uuid.New(),
		Name:        name,
		Description: description,
	}, nil
}

func (s *Manager) CreateGroup(ctx Context, name, description string) (*Group, error) {
	if err := s.validateName(name); err != nil {
		return nil, fmt.Errorf("CreateGroup: %s", err.Error())
	}

	return &Group{
		GroupID:     uuid.New(),
		Name:        name,
		Description: description,
	}, nil
}

func (s *Manager) CreatePermission(ctx Context, name, description string) (*Permission, error) {
	if err := s.validateName(name); err != nil {
		return nil, fmt.Errorf("CreatePermission: %s", err.Error())
	}

	return &Permission{
		PermissionID: uuid.New(),
		Name:         name,
		Description:  description,
	}, nil
}

func (s *Manager) SetupPrincipalAsOwner(ctx Context, pri *Principal) (*Group, *Role, *Permission, error) {
	g, err := pri.AttachNewGroup(GroupOwnerName, GroupOwnerDescription)
	if err != nil {
		return nil, nil, nil, err
	}

	r, err := g.AttachNewRole(RoleOwnerName, RoleOwnerDescription)
	if err != nil {
		return nil, nil, nil, err
	}

	perm, err := r.AttachNewPermission(ctx, PermissionOwnerName, PermissionOwnerDescription)
	if err != nil {
		return nil, nil, nil, err
	}

	return g, r, perm, nil

}

func (s *Manager) EnumerateGroupIDs(ctx Context) ([]string, error) {
	IDs, err := s.q.EnumerateGroupIDs(ctx)
	if err != nil {
		return nil, err
	}

	if len(IDs) == 0 {
		return []string{}, nil
	}

	result := make([]string, 0, len(IDs))
	for _, u := range IDs {
		result = append(result, u.String())
	}

	return result, nil
}

func (s *Manager) FindGroups(ctx Context, ids []string) ([]*Group, error) {
	return s.q.FindGroups(ctx, ids)
}

func (s *Manager) FindGroupByID(ctx Context, groupID string) (*Group, error) {
	g, err := s.q.FindGroups(ctx, []string{groupID})
	if err != nil {
		return nil, err
	}
	if len(g) == 0 {
		log.Println("FindGroupByID:", groupID, "no entry")
		return nil, ErrNoEntry
	}

	return g[0], nil
}

func (s *Manager) EnumerateRoleIDs(ctx Context) ([]string, error) {
	IDs, err := s.q.EnumerateRoleIDs(ctx)
	if err != nil {
		return nil, err
	}

	if len(IDs) == 0 {
		return []string{}, nil
	}

	result := make([]string, 0, len(IDs))
	for _, u := range IDs {
		result = append(result, u.String())
	}

	return result, nil
}

func (s *Manager) FindRoles(ctx Context, ids []string) ([]*Role, error) {
	return s.q.FindRoles(ctx, ids)
}

func (s *Manager) FindRoleByID(ctx Context, roleID string) (*Role, error) {
	r, err := s.q.FindRoles(ctx, []string{roleID})
	if err != nil {
		return nil, err
	}
	if len(r) == 0 {
		log.Println("FindRoleByID:", roleID, "no entry")
		return nil, ErrNoEntry
	}

	return r[0], nil
}

func (s *Manager) EnumeratePermissionIDs(ctx Context) ([]string, error) {
	IDs, err := s.q.EnumeratePermissionIDs(ctx)
	if err != nil {
		return nil, err
	}

	if len(IDs) == 0 {
		return []string{}, nil
	}

	result := make([]string, 0, len(IDs))
	for _, u := range IDs {
		result = append(result, u.String())
	}

	return result, nil
}

func (s *Manager) FindPermissions(ctx Context, ids []string) ([]*Permission, error) {
	return s.q.FindPermissions(ctx, ids)
}

func (s *Manager) FindPermissionByID(ctx Context, permissionID string) (*Permission, error) {
	r, err := s.q.FindPermissions(ctx, []string{permissionID})
	if err != nil {
		return nil, err
	}
	if len(r) == 0 {
		log.Println("FindPermissionByID:", permissionID, "no entry")
		return nil, ErrNoEntry
	}

	return r[0], nil
}

func (s *Manager) CreateMappingRule(ctx Context, name, description string, ruleType MappingType, priority int,
	detail, associationType, associationID string) (*MappingRule, error) {

	if err := s.validateName(name); err != nil {
		return nil, fmt.Errorf("CreateMappingRule: %s", err.Error())
	}

	var group *Group
	var role *Role
	switch associationType {
	case "group":
		groups, err := s.FindGroups(ctx, []string{associationID})
		if err != nil || len(groups) == 0 {
			return nil, fmt.Errorf("CreateMappingRule: specified group is not found")
		}
		group = groups[0]

	case "role":
		roles, err := s.FindRoles(ctx, []string{associationID})
		if err != nil || len(roles) == 0 {
			return nil, fmt.Errorf("CreateMappingRule: specified role is not found")
		}
		role = roles[0]
	}

	rule := &MappingRule{
		MappingRuleID:   uuid.New(),
		RuleType:        MappingType(ruleType),
		Detail:          detail,
		Name:            name,
		Description:     description,
		Priority:        priority, // TODO
		associatedGroup: group,
		associatedRole:  role,
	}

	return rule, nil
}

func (s *Manager) EnumerateMappingRuleIDs(ctx Context) ([]string, error) {
	IDs, err := s.q.EnumerateMappingRuleIDs(ctx)
	if err != nil {
		return nil, err
	}

	if len(IDs) == 0 {
		return []string{}, nil
	}

	result := make([]string, 0, len(IDs))
	for _, u := range IDs {
		result = append(result, u.String())
	}

	return result, nil
}

func (s *Manager) FindMappingRules(ctx Context, ids []string) ([]*MappingRule, error) {
	return s.q.FindMappingRules(ctx, ids)
}

func (s *Manager) EnumerateMappingRules(ctx Context) ([]*MappingRule, error) {
	IDs, err := s.q.EnumerateMappingRuleIDs(ctx)
	if err != nil {
		return nil, err
	}

	if len(IDs) == 0 {
		return []*MappingRule{}, nil
	}

	result := make([]string, 0, len(IDs))
	for _, u := range IDs {
		result = append(result, u.String())
	}

	rules, err := s.q.FindMappingRules(ctx, result)
	if err != nil {
		return nil, err
	}

	return rules, nil
}
