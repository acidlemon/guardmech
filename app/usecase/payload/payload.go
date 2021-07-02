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
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
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
		Principal: &Principal{
			ID:          pri.PrincipalID,
			Name:        pri.Name,
			Description: pri.Description,
		},
		Auth: &Auth{
			ID:      oidcAuth.OIDCAuthID,
			Issuer:  oidcAuth.Issuer,
			Subject: oidcAuth.Subject,
			Email:   oidcAuth.Email,
		},
		APIKeys:     make([]*APIKey, 0, len(apikeys)),
		Groups:      make([]*Group, 0, len(gs)),
		Roles:       make([]*Role, 0, len(rs)),
		Permissions: make([]*Permission, 0, len(perms)),
	}

	for _, v := range apikeys {
		result.APIKeys = append(result.APIKeys, &APIKey{
			ID:          v.AuthAPIKeyID,
			Name:        v.Name,
			MaskedToken: v.MaskedToken,
		})
	}
	for _, v := range gs {
		result.Groups = append(result.Groups, &Group{
			ID:          v.GroupID,
			Name:        v.Name,
			Description: v.Description,
		})
	}
	for _, v := range rs {
		result.Roles = append(result.Roles, &Role{
			ID:          v.RoleID,
			Name:        v.Name,
			Description: v.Description,
		})
	}
	for _, v := range perms {
		result.Permissions = append(result.Permissions, &Permission{
			ID:          v.PermissionID,
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return result
}
