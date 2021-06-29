package membership

import (
	"github.com/google/uuid"
)

type factory struct {
	unused Query
}

type Factory interface {
	NewPrincipal(
		ID uuid.UUID,
		name, description string,
		auth *OIDCAuthorization,
		apikeys []*AuthAPIKey,
		roles []*Role,
		groups []*Group,
	) *Principal

	NewRole(
		ID uuid.UUID,
		name, description string,
		perms []*Permission,
	) *Role

	NewGroup(
		ID uuid.UUID,
		name, description string,
		roles []*Role,
	) *Group
}

func NewFactory(q Query) Factory {
	return &factory{unused: q}
}

func (f *factory) NewPrincipal(
	ID uuid.UUID,
	name, description string,
	auth *OIDCAuthorization,
	apikeys []*AuthAPIKey,
	roles []*Role,
	groups []*Group,
) *Principal {
	return &Principal{
		PrincipalID: ID,
		Name:        name,
		Description: description,

		auth:    auth,
		apiKeys: apikeys,
		roles:   roles,
		groups:  groups,
	}
}

func (f *factory) NewRole(
	ID uuid.UUID,
	name, description string,
	perms []*Permission,
) *Role {
	return &Role{
		RoleID:      ID,
		Name:        name,
		Description: description,
		permissions: perms,
	}
}

func (f *factory) NewGroup(
	ID uuid.UUID,
	name, description string,
	roles []*Role,
) *Group {
	return &Group{
		GroupID:     ID,
		Name:        name,
		Description: description,
		roles:       roles,
	}
}
