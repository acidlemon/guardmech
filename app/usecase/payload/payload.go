package payload

import (
	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/google/uuid"
)

type Auth struct {
	ID        uuid.UUID  `json:"id"`
	Issuer    string     `json:"issuer"`
	Subject   string     `json:"subject"`
	Email     string     `json:"email"`
	Principal *Principal `json:"-" db:"-"`
}

type Principal struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type APIKey struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	MaskedToken string    `json:"masked_token"`
}

type Group struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	AttachedRoles []*Role   `json:"attached_roles"`
}

type Role struct {
	ID                  uuid.UUID     `json:"id"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	AttachedPermissions []*Permission `json:"attached_permissions"`
}

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type MappingRule struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	RuleType        int       `json:"rule_type"`
	Detail          string    `json:"detail"`
	Priority        int       `json:"priority"`
	AssociationType string    `json:"association_type"`
	AssociationID   string    `json:"association_id"`
}

type PrincipalPayload struct {
	Principal   *Principal    `json:"principal"`
	Auth        *Auth         `json:"auth_oidc"`
	APIKeys     []*APIKey     `json:"auth_apikeys"`
	Groups      []*Group      `json:"groups"`
	Roles       []*Role       `json:"roles"`
	Permissions []*Permission `json:"permissions"`
}

func PrincipalFromEntity(pri *entity.Principal) *Principal {
	return &Principal{
		ID:          pri.PrincipalID,
		Name:        pri.Name,
		Description: pri.Description,
	}
}

func PrincipalPayloadFromEntity(pri *entity.Principal) *PrincipalPayload {
	oidcAuth := pri.OIDCAuthorization()
	apikeys := pri.APIKeys()
	gs := pri.Groups()
	rs := pri.Roles()
	perms := pri.Permissions()

	result := &PrincipalPayload{
		Principal:   PrincipalFromEntity(pri),
		Auth:        AuthFromEntity(oidcAuth),
		APIKeys:     make([]*APIKey, 0, len(apikeys)),
		Groups:      make([]*Group, 0, len(gs)),
		Roles:       make([]*Role, 0, len(rs)),
		Permissions: make([]*Permission, 0, len(perms)),
	}

	for _, v := range apikeys {
		result.APIKeys = append(result.APIKeys, APIKeyFromEntity(v))
	}
	for _, v := range gs {
		result.Groups = append(result.Groups, GroupFromEntity(v))
	}
	for _, v := range rs {
		result.Roles = append(result.Roles, RoleFromEntity(v))
	}
	for _, v := range perms {
		result.Permissions = append(result.Permissions, PermissionFromEntity(v))
	}

	return result
}

func AuthFromEntity(a *entity.OIDCAuthorization) *Auth {
	if a == nil {
		return nil
	}

	return &Auth{
		ID:      a.OIDCAuthID,
		Issuer:  a.Issuer,
		Subject: a.Subject,
		Email:   a.Email,
	}
}

func APIKeyFromEntity(a *entity.AuthAPIKey) *APIKey {
	return &APIKey{
		ID:          a.AuthAPIKeyID,
		Name:        a.Name,
		MaskedToken: a.MaskedToken,
	}
}

func RoleFromEntity(r *entity.Role) *Role {
	perms := r.Permissions()
	attachedPerms := make([]*Permission, 0, len(perms))
	for _, v := range perms {
		attachedPerms = append(attachedPerms, PermissionFromEntity(v))
	}

	return &Role{
		ID:                  r.RoleID,
		Name:                r.Name,
		Description:         r.Description,
		AttachedPermissions: attachedPerms,
	}
}

func GroupFromEntity(g *entity.Group) *Group {
	roles := g.Roles()
	attachedRoles := make([]*Role, 0, len(roles))
	for _, v := range roles {
		attachedRoles = append(attachedRoles, RoleFromEntity(v))
	}

	return &Group{
		ID:            g.GroupID,
		Name:          g.Name,
		Description:   g.Description,
		AttachedRoles: attachedRoles,
	}
}

func PermissionFromEntity(perm *entity.Permission) *Permission {
	return &Permission{
		ID:          perm.PermissionID,
		Name:        perm.Name,
		Description: perm.Description,
	}
}

func MappingRuleFromEntity(rule *entity.MappingRule) *MappingRule {
	associationType := ""
	associationID := ""
	if group := rule.AssociatedGroup(); group != nil {
		associationType = "group"
		associationID = group.GroupID.String()
	} else if role := rule.AssociatedRole(); role != nil {
		associationType = "role"
		associationID = role.RoleID.String()
	}

	return &MappingRule{
		ID:              rule.MappingRuleID,
		Name:            rule.Name,
		Description:     rule.Description,
		Priority:        rule.Priority,
		RuleType:        int(rule.RuleType),
		Detail:          rule.Detail,
		AssociationType: associationType,
		AssociationID:   associationID,
	}
}
