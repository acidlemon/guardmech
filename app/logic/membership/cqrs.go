package membership

import (
	"errors"

	"github.com/google/uuid"
)

var ErrNoEntry = errors.New("no such entry")

type Query interface {
	FindPrincipals(ctx Context, ids []string) ([]*Principal, error)
	FindPrincipalIDByOIDC(ctx Context, issuer, subject string) (*Principal, error)
	EnumeratePrincipalIDs(ctx Context) ([]uuid.UUID, error)

	// FindGroups(ctx Context, ids []string) ([]*Group, error)
	// EnumerateGroupIDs(ctx Context) ([]uuid.UUID, error)

	FindRoles(ctx Context, ids []string) ([]*Role, error)
	EnumerateRoleIDs(ctx Context) ([]uuid.UUID, error)

	FindPermissions(ctx Context, ids []string) ([]*Permission, error)
	EnumeratePermissionIDs(ctx Context) ([]uuid.UUID, error)
}

type Command interface {
	SavePrincipal(ctx Context, pri *Principal) error
	SaveRole(ctx Context, r *Role) error
	SavePermission(ctx Context, perm *Permission) error
	SaveAuthOIDC(ctx Context, oidc *OIDCAuthorization, pri *Principal) error
	SaveAuthAPIKey(ctx Context, key *AuthAPIKey, pri *Principal) error
}
