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

	FindGroups(ctx Context, ids []string) ([]*Group, error)
	EnumerateGroupIDs(ctx Context) ([]uuid.UUID, error)

	FindRoles(ctx Context, ids []string) ([]*Role, error)
	EnumerateRoleIDs(ctx Context) ([]uuid.UUID, error)

	FindPermissions(ctx Context, ids []string) ([]*Permission, error)
	EnumeratePermissionIDs(ctx Context) ([]uuid.UUID, error)

	FindMappingRules(ctx Context, ids []string) ([]*MappingRule, error)
	EnumerateMappingRuleIDs(ctx Context) ([]uuid.UUID, error)
}

type Command interface {
	Error() error // see https://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike

	SavePrincipal(ctx Context, pri *Principal)
	SaveGroup(ctx Context, g *Group)
	SaveRole(ctx Context, r *Role)
	SavePermission(ctx Context, perm *Permission)
	SaveAuthOIDC(ctx Context, oidc *OIDCAuthorization, pri *Principal)
	SaveAuthAPIKey(ctx Context, key *AuthAPIKey, pri *Principal)

	DeletePrincipal(ctx Context, pri *Principal)
	DeleteGroup(ctx Context, g *Group)
	DeleteRole(ctx Context, r *Role)
	DeletePermission(ctx Context, perm *Permission)
}
