package persistence

import (
	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/guardmech/persistence/db"
	"github.com/acidlemon/seacle"
)

type command struct {
	conn seacle.Executable
	m    *db.Service
}

func NewCommand(conn seacle.Executable) entity.Command {
	return &command{
		conn: conn,
		m:    &db.Service{},
	}
}

func (c *command) SavePrincipal(ctx Context, pri *entity.Principal) error {
	return c.m.SavePrincipal(ctx, c.conn, pri)
}

func (c *command) SaveRole(ctx Context, r *entity.Role) error {
	return c.m.SaveRole(ctx, c.conn, r)
}

func (c *command) SavePermission(ctx Context, perm *entity.Permission) error {
	return c.m.SavePermission(ctx, c.conn, perm)
}

func (c *command) SaveAuthOIDC(ctx Context, oidc *entity.OIDCAuthorization, pri *entity.Principal) error {
	return c.m.SaveAuthOIDC(ctx, c.conn, oidc, pri)
}

func (c *command) SaveAuthAPIKey(ctx Context, key *entity.AuthAPIKey, pri *entity.Principal) error {
	return c.m.SaveAuthAPIKey(ctx, c.conn, key, pri)
}
