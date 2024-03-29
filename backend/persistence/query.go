package persistence

import (
	"context"

	entity "github.com/acidlemon/guardmech/backend/app/logic/membership"
	"github.com/acidlemon/guardmech/backend/persistence/db"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type query struct {
	conn seacle.Selectable
	m    *db.Service
}

func NewQuery(conn seacle.Selectable) entity.Query {
	return &query{
		conn: conn,
		m:    &db.Service{},
	}
}

type Context = context.Context

func (q *query) FindPrincipals(ctx Context, ids []string) ([]*entity.Principal, error) {
	return q.m.FindPrincipals(ctx, q.conn, ids)
}

func (q *query) FindPrincipalIDByOIDC(ctx Context, issuer, subject string) (*entity.Principal, error) {
	return q.m.FindPrincipalByOIDC(ctx, q.conn, issuer, subject)
}

func (q *query) EnumeratePrincipalIDs(ctx Context) ([]uuid.UUID, error) {
	ids, err := q.m.EnumeratePrincipalIDs(ctx, q.conn)
	if err != nil {
		return nil, err
	}

	result := make([]uuid.UUID, 0, len(ids))
	for _, v := range ids {
		result = append(result, uuid.MustParse(v))
	}

	return result, nil
}

func (q *query) FindGroups(ctx Context, ids []string) ([]*entity.Group, error) {
	return q.m.FindGroups(ctx, q.conn, ids)
}

func (q *query) EnumerateGroupIDs(ctx Context) ([]uuid.UUID, error) {
	ids, err := q.m.EnumerateGroupIDs(ctx, q.conn)
	if err != nil {
		return nil, err
	}

	result := make([]uuid.UUID, 0, len(ids))
	for _, v := range ids {
		result = append(result, uuid.MustParse(v))
	}

	return result, nil
}

func (q *query) FindRoles(ctx Context, ids []string) ([]*entity.Role, error) {
	return q.m.FindRoles(ctx, q.conn, ids)
}

func (q *query) EnumerateRoleIDs(ctx Context) ([]uuid.UUID, error) {
	ids, err := q.m.EnumerateRoleIDs(ctx, q.conn)
	if err != nil {
		return nil, err
	}

	result := make([]uuid.UUID, 0, len(ids))
	for _, v := range ids {
		result = append(result, uuid.MustParse(v))
	}

	return result, nil
}

func (q *query) FindPermissions(ctx Context, ids []string) ([]*entity.Permission, error) {
	return q.m.FindPermissions(ctx, q.conn, ids)
}

func (q *query) EnumeratePermissionIDs(ctx Context) ([]uuid.UUID, error) {
	ids, err := q.m.EnumeratePermissionIDs(ctx, q.conn)
	if err != nil {
		return nil, err
	}

	result := make([]uuid.UUID, 0, len(ids))
	for _, v := range ids {
		result = append(result, uuid.MustParse(v))
	}

	return result, nil
}

func (q *query) FindMappingRules(ctx Context, ids []string) ([]*entity.MappingRule, error) {
	return q.m.FindMappingRules(ctx, q.conn, ids)
}

func (q *query) EnumerateMappingRuleIDs(ctx Context) ([]uuid.UUID, error) {
	ids, err := q.m.EnumerateMappingRuleIDs(ctx, q.conn)
	if err != nil {
		return nil, err
	}

	result := make([]uuid.UUID, 0, len(ids))
	for _, v := range ids {
		result = append(result, uuid.MustParse(v))
	}

	return result, nil
}
