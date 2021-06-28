package persistence

import (
	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/guardmech/persistence/db"
	"github.com/acidlemon/seacle"
)

type command struct {
	conn seacle.Executable
}

func NewCommand(conn seacle.Executable) entity.Command {
	return &command{
		conn: conn,
	}
}

func (c *command) SavePrincipal(ctx Context, pri *entity.Principal) error {
	m := db.Service{}

	return m.SavePrincipal(ctx, c.conn, pri)
}

func (c *command) SaveRole(ctx Context, r *entity.Role) error {
	m := db.Service{}

	return m.SaveRole(ctx, c.conn, r)
}

func (c *command) SavePermission(ctx Context, perm *entity.Permission) error {
	m := db.Service{}

	return m.SavePermission(ctx, c.conn, perm)
}

func (c *command) SaveAuthOIDC(ctx Context, oidc *entity.OIDCAuthorization, pri *entity.Principal) error {
	m := db.Service{}

	return m.SaveAuthOIDC(ctx, c.conn, oidc, pri)
}
