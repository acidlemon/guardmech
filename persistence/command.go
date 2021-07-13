package persistence

import (
	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/guardmech/persistence/db"
	"github.com/acidlemon/seacle"
)

type command struct {
	conn seacle.Executable
	m    *db.Service
	err  error
}

func NewCommand(conn seacle.Executable) entity.Command {
	return &command{
		conn: conn,
		m:    &db.Service{},
		err:  nil,
	}
}

func (c *command) Error() error {
	err := c.err
	c.err = nil // clear error state
	return err
}

func (c *command) SavePrincipal(ctx Context, pri *entity.Principal) {
	if c.err != nil {
		return
	}
	c.err = c.m.SavePrincipal(ctx, c.conn, pri)
}

func (c *command) SaveGroup(ctx Context, g *entity.Group) {
	if c.err != nil {
		return
	}
	c.err = c.m.SaveGroup(ctx, c.conn, g)
}

func (c *command) SaveRole(ctx Context, r *entity.Role) {
	if c.err != nil {
		return
	}
	c.err = c.m.SaveRole(ctx, c.conn, r)
}

func (c *command) SavePermission(ctx Context, perm *entity.Permission) {
	if c.err != nil {
		return
	}
	c.err = c.m.SavePermission(ctx, c.conn, perm)
}

func (c *command) SaveAuthOIDC(ctx Context, oidc *entity.OIDCAuthorization, pri *entity.Principal) {
	if c.err != nil {
		return
	}
	c.err = c.m.SaveAuthOIDC(ctx, c.conn, oidc, pri)
}

func (c *command) SaveAuthAPIKey(ctx Context, key *entity.AuthAPIKey, pri *entity.Principal) {
	if c.err != nil {
		return
	}
	c.err = c.m.SaveAuthAPIKey(ctx, c.conn, key, pri)
}

func (c *command) DeletePrincipal(ctx Context, pri *entity.Principal) {
	if c.err != nil {
		return
	}
	c.err = c.m.DeletePrincipal(ctx, c.conn, pri)
}

func (c *command) DeleteGroup(ctx Context, g *entity.Group) {
	if c.err != nil {
		return
	}
	c.err = c.m.DeleteGroup(ctx, c.conn, g)
}

func (c *command) DeleteRole(ctx Context, r *entity.Role) {
	if c.err != nil {
		return
	}
	c.err = c.m.DeleteRole(ctx, c.conn, r)
}

func (c *command) DeletePermission(ctx Context, perm *entity.Permission) {
	if c.err != nil {
		return
	}
	c.err = c.m.DeletePermission(ctx, c.conn, perm)
}
