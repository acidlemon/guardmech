package payload

import (
	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/google/uuid"
)

type Auth struct {
	UniqueID  uuid.UUID  `json:"unique_id"`
	Issuer    string     `json:"issuer"`
	Subject   string     `json:"subject"`
	Email     string     `json:"email"`
	Principal *Principal `json:"-" db:"-"`
}

type Principal struct {
	UniqueID    uuid.UUID `json:"unique_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type APIKey struct {
}

type Group struct {
	UniqueID    uuid.UUID `json:"unique_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Role struct {
	UniqueID    uuid.UUID `json:"unique_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Permission struct {
	UniqueID    uuid.UUID `json:"unique_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type PrincipalPayload struct {
	Principal   *Principal    `json:"principal"`
	Auth        *Auth         `json:"auth_oidc"`
	APIKeys     []*APIKey     `json:"auth_api_keys"`
	Groups      []*Group      `json:"groups"`
	Roles       []*Role       `json:"roles"`
	Permissions []*Permission `json:"permissions"`
}

func PrincipalPayloadFromEntity(pri *entity.Principal) *PrincipalPayload {
	gs := pri.Groups()
	rs := pri.Roles()
	perms := pri.Permissions()

	result := &PrincipalPayload{
		APIKeys:     []*APIKey{},
		Groups:      make([]*Group, 0, len(gs)),
		Roles:       make([]*Role, 0, len(rs)),
		Permissions: make([]*Permission, 0, len(perms)),
	}

	result.Principal = &Principal{
		UniqueID:    pri.PrincipalID,
		Name:        pri.Name,
		Description: pri.Description,
	}
	result.Auth = &Auth{}

	for _, v := range gs {
		result.Groups = append(result.Groups, &Group{
			UniqueID:    v.GroupID,
			Name:        v.Name,
			Description: v.Description,
		})
	}
	for _, v := range rs {
		result.Roles = append(result.Roles, &Role{
			UniqueID:    v.RoleID,
			Name:        v.Name,
			Description: v.Description,
		})
	}
	for _, v := range perms {
		result.Permissions = append(result.Permissions, &Permission{
			UniqueID:    v.PermissionID,
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return result
}
