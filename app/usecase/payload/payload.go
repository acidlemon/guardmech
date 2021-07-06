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
	return &Role{
		ID:          r.RoleID,
		Name:        r.Name,
		Description: r.Description,
	}
}

func GroupFromEntity(g *entity.Group) *Group {
	return &Group{
		ID:          g.GroupID,
		Name:        g.Name,
		Description: g.Description,
	}
}

func PermissionFromEntity(perm *entity.Permission) *Permission {
	return &Permission{
		ID:          perm.PermissionID,
		Name:        perm.Name,
		Description: perm.Description,
	}
}
