package db

import (
	"database/sql"
	"fmt"
	"log"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type PrincipalRow struct {
	SeqID       int64  `db:"seq_id,primary"`
	PrincipalID string `db:"principal_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func principalRowFromEntity(pri *entity.Principal) *PrincipalRow {
	return &PrincipalRow{
		PrincipalID: pri.PrincipalID.String(),
		Name:        pri.Name,
		Description: pri.Description,
	}
}
func (pri *PrincipalRow) ToEntity() *entity.Principal {
	return &entity.Principal{
		PrincipalID: uuid.MustParse(pri.PrincipalID),
		Name:        pri.Name,
		Description: pri.Description,
	}
}

type Service struct {
}

// func (s *Service) HasPrincipal(ctx context.Context, conn *sql.Conn) (bool, error) {
// 	row := conn.QueryRowContext(ctx, `SELECT COUNT(*) AS cnt FROM principal`)
// 	var cnt int
// 	err := row.Scan(&cnt)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	if cnt == 0 {
// 		return false, nil
// 	}
// 	return true, nil
// }

func (s *Service) FindPrincipalBySeqID(ctx Context, conn *sql.Conn, id int64) (*entity.Principal, error) {
	pr := &PrincipalRow{}
	err := seacle.SelectRow(ctx, conn, pr, `WHERE seq_id = ?`, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	mempr := pr.ToEntity()

	return mempr, nil
}

/*
func (s *Membership) FetchPrincipalPayload(ctx Context, conn *sql.Conn, id int64) (*membership.PrincipalPayload, error) {
	log.Println(string(debug.Stack()))

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
	memauths := make([]*entity.Auth, 0, len(auths))
	for _, v := range auths {
		a := membership.Auth(v.Auth)
		a.Principal = &mempr
		memauths = append(memauths, &a)
	}
	memapikeys := make([]*entity.APIKey, 0, len(apikeys))
	for _, v := range apikeys {
		a := membership.APIKey(v.APIKey)
		a.Principal = &mempr
		memapikeys = append(memapikeys, &a)
	}
	memgroups := make([]*entity.Group, 0, len(groups))
	for _, v := range groups {
		g := membership.Group(*v)
		memgroups = append(memgroups, &g)
	}
	memroles := make([]*entity.Role, 0, len(roles))
	for _, v := range roles {
		r := membership.Role(*v)
		memroles = append(memroles, &r)
	}
	memperms := make([]*entity.Permission, 0, len(perms))
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
*/

func (s *Service) SavePrincipal(ctx Context, conn seacle.Executable, pri *entity.Principal) error {
	row := &PrincipalRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE principal_id = ?", pri.PrincipalID.String())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return s.createPrincipal(ctx, conn, pri)
	} else {
		return s.updatePrincipal(ctx, conn, pri, row)
	}
}

func (s *Service) createPrincipal(ctx Context, conn seacle.Executable, pri *entity.Principal) error {
	row := principalRowFromEntity(pri)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updatePrincipal(ctx Context, conn seacle.Executable, pri *entity.Principal, row *PrincipalRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `FOR UPDATE WHERE seq_id = ?`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock principal row: err=%s", err)
	}

	// update row
	row.Name = pri.Name
	row.Description = pri.Description
	err = seacle.Update(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update principal row: err=%s", err)
	}

	return nil
}

func (s *Service) FindPrincipals(ctx Context, conn seacle.Selectable, principalIDs []string) ([]*entity.Principal, error) {
	pris := make([]*PrincipalRow, 0, 8)
	err := seacle.Select(ctx, conn, &pris, `WHERE principal_id IN (?)`, principalIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]*entity.Principal, 0, len(pris))

	// seq_idを抽出
	priSeqIDs := make([]int64, 0, len(pris))
	for _, v := range pris {
		result = append(result, v.ToEntity())
		priSeqIDs = append(priSeqIDs, v.SeqID)
	}

	// Role
	// rs := make([]*RoleRow, 0, 12)
	// err = seacle.Select(ctx, conn, &rs, `JOIN `)

	return result, nil
}

func (s *Service) FindPrincipalByOIDC(ctx Context, conn seacle.Selectable, issuer, subject string) (*entity.Principal, error) {
	pri := PrincipalRow{}

	log.Println("issuer=", issuer, "subject=", subject)

	err := seacle.SelectRow(ctx, conn, &pri,
		`JOIN auth_oidc a ON a.principal_seq_id = principal.seq_id WHERE a.issuer = ? AND a.subject = ?`,
		issuer, subject)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := s.FindPrincipals(ctx, conn, []string{pri.PrincipalID})
	if err != nil {
		return nil, err
	}

	return result[0], nil
}

func (s *Service) EnumeratePrincipalIDs(ctx Context, conn seacle.Selectable) ([]string, error) {
	pris := make([]*PrincipalRow, 0, 8)
	err := seacle.Select(ctx, conn, &pris, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]string, 0, len(pris))
	for _, v := range pris {
		result = append(result, v.PrincipalID)
	}

	return result, nil
}

// func (s *Service) FetchAllPrincipal(ctx Context, conn seacle.Selectable) ([]*membership.Principal, error) {
// 	principals := make([]*Principal, 0, 4)
// 	err := seacle.Select(ctx, conn, &principals, ``)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result := []*membership.Principal{}
// 	for _, v := range principals {
// 		pr := membership.Principal(*v)
// 		result = append(result, &pr)
// 	}

// 	return result, nil
// }
