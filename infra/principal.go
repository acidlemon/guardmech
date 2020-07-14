package infra

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/acidlemon/guardmech/db"
	"github.com/acidlemon/guardmech/membership"
	"github.com/acidlemon/seacle"
)

type Principal membership.Principal

type Membership struct {
}

func (s *Membership) HasPrincipal(ctx context.Context, conn *sql.Conn) (bool, error) {
	row := conn.QueryRowContext(ctx, `SELECT COUNT(*) AS cnt FROM principal`)
	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return false, err
	}

	if cnt == 0 {
		return false, nil
	}
	return true, nil
}

func (s *Membership) FetchPrincipalPayload(ctx context.Context, conn *sql.Conn, id int64) (*membership.PrincipalPayload, error) {
	pr := &Principal{}
	err := seacle.SelectRow(ctx, conn, pr, `WHERE seq_id = ?`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	mempr := membership.Principal(*pr)

	auths := make([]*Auth, 0, 4)
	err = seacle.Select(ctx, conn, &auths, `WHERE principal_id = ?`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	apikeys := make([]*APIKey, 0, 4)
	err = seacle.Select(ctx, conn, &apikeys, `WHERE principal_id = ?`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	groups := make([]*Group, 0, 8)
	err = seacle.Select(ctx, conn, &groups,
		`JOIN principal_group_map AS m ON group_info.seq_id = m.group_id WHERE m.principal_id = ?`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	roles := make([]*Role, 0, 8)
	err = seacle.Select(ctx, conn, &roles,
		`JOIN principal_role_map AS m ON role_info.seq_id = m.role_id WHERE m.principal_id = ?`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(groups) > 0 {
		groupIDs := make([]int64, 0, len(groups))
		for _, v := range groups {
			groupIDs = append(groupIDs, v.SeqID)
		}

		// append to roles
		err = seacle.Select(ctx, conn, &roles,
			`JOIN group_role_map AS m ON role_info.seq_id = m.role_id WHERE m.group_id IN (?)`, groupIDs)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		roles = uniqRoles(roles)
	}

	perms := make([]*Permission, 0, 8)
	if len(roles) > 0 {
		err = seacle.Select(ctx, conn, &perms,
			`JOIN role_permission_map AS m ON permission.seq_id = m.permission_id WHERE m.role_id IN (?)`, roleIDs(roles))
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	// 詰め替え
	memauths := make([]*membership.Auth, 0, len(auths))
	for _, v := range auths {
		a := membership.Auth(v.Auth)
		a.Principal = &mempr
		memauths = append(memauths, &a)
	}
	memapikeys := make([]*membership.APIKey, 0, len(apikeys))
	for _, v := range apikeys {
		a := membership.APIKey(v.APIKey)
		a.Principal = &mempr
		memapikeys = append(memapikeys, &a)
	}
	memgroups := make([]*membership.Group, 0, len(groups))
	for _, v := range groups {
		g := membership.Group(*v)
		memgroups = append(memgroups, &g)
	}
	memroles := make([]*membership.Role, 0, len(roles))
	for _, v := range roles {
		r := membership.Role(*v)
		memroles = append(memroles, &r)
	}
	memperms := make([]*membership.Permission, 0, len(perms))
	for _, v := range perms {
		pe := membership.Permission(*v)
		memperms = append(memperms, &pe)
	}

	return &membership.PrincipalPayload{
		Principal:   &mempr,
		Auths:       memauths,
		APIKeys:     memapikeys,
		Groups:      memgroups,
		Roles:       memroles,
		Permissions: memperms,
	}, nil
}

func (s *Membership) SavePrincipal(ctx context.Context, tx *db.Tx, pri *membership.Principal) (int64, error) {
	if pri.SeqID == 0 {
		return s.createPrincipal(ctx, tx, pri)
	} else {
		return s.updatePrincipal(ctx, tx, pri)
	}
}

func (s *Membership) createPrincipal(ctx context.Context, tx *db.Tx, pri *membership.Principal) (int64, error) {
	// new Principal
	principal := Principal(*pri)
	seqID, err := seacle.Insert(ctx, tx, &principal)
	if err != nil {
		return 0, err
	}

	return seqID, nil
}

func (s *Membership) updatePrincipal(ctx context.Context, tx *db.Tx, pri *membership.Principal) (int64, error) {
	priForUpdate := &Principal{}
	err := seacle.SelectRow(ctx, tx, priForUpdate, `WHERE unique_id = ?`, pri.UniqueID)
	if err != nil {
		// TODO: fallback to createRole
		return 0, nil
	}

	if priForUpdate.SeqID != pri.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, r.SeqID=%d", priForUpdate.SeqID, pri.SeqID)
	}

	// lock row
	err = seacle.SelectRow(ctx, tx, priForUpdate, `FOR UPDATE WHERE seq_id = ?`, priForUpdate.SeqID)
	if err != nil {
		return 0, fmt.Errorf("failed to lock principal row: err=%s", err)
	}

	// update row
	principal := Principal(*pri)
	err = seacle.Update(ctx, tx, &principal)
	if err != nil {
		return 0, fmt.Errorf("failed to update principal row: err=%s", err)
	}

	return pri.SeqID, nil
}

func (s *Membership) FindPrincipal(ctx context.Context, conn *sql.Conn, issuer, subject string) (*membership.Principal, error) {
	pr := &Principal{}
	err := seacle.SelectRow(ctx, conn, pr,
		`JOIN auth AS a ON a.principal_id = principal.seq_id WHERE a.issuer = ? AND a.subject = ?`,
		issuer, subject)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	mempr := membership.Principal(*pr)

	return &mempr, nil
}

func (s *Membership) FetchAllPrincipal(ctx context.Context, conn *sql.Conn) ([]*membership.Principal, error) {
	principals := make([]*Principal, 0, 4)
	err := seacle.Select(ctx, conn, &principals, ``)
	if err != nil {
		return nil, err
	}

	result := []*membership.Principal{}
	for _, v := range principals {
		pr := membership.Principal(*v)
		result = append(result, &pr)
	}

	return result, nil
}
