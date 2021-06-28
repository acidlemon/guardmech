package persistence

import (
	"context"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/guardmech/persistence/db"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type query struct {
	conn seacle.Selectable
}

func NewQuery(conn seacle.Selectable) entity.Query {
	return &query{
		conn: conn,
	}
}

type Context = context.Context

func (q *query) FindPrincipals(ctx Context, ids []string) ([]*entity.Principal, error) {
	m := db.Service{}
	result, err := m.FindPrincipals(ctx, q.conn, ids)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (q *query) FindPrincipalIDByOIDC(ctx Context, issuer, subject string) (*entity.Principal, error) {
	m := db.Service{}
	pri, err := m.FindPrincipalByOIDC(ctx, q.conn, issuer, subject)
	if err != nil {
		return nil, err
	}

	return pri, nil
}

func (q *query) EnumeratePrincipalIDs(ctx Context) ([]uuid.UUID, error) {
	m := db.Service{}
	ids, err := m.EnumeratePrincipalIDs(ctx, q.conn)
	if err != nil {
		return nil, err
	}

	result := make([]uuid.UUID, 0, len(ids))
	for _, v := range ids {
		result = append(result, uuid.MustParse(v))
	}

	return result, nil
}
