package db

import (
	"database/sql"
	"fmt"
	"log"

	entity "github.com/acidlemon/guardmech/backend/app/logic/membership"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type PrincipalRow struct {
	SeqID       int64  `db:"seq_id,primary,auto_increment"`
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
	q entity.Query
}

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

func (s *Service) SavePrincipal(ctx Context, conn seacle.Executable, pri *entity.Principal) error {
	row := &PrincipalRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE principal_id = ?", pri.PrincipalID.String())
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	if err == sql.ErrNoRows {
		err = s.createPrincipalRow(ctx, conn, pri)
		if err != nil {
			log.Println(err)
			return err
		}
		err = seacle.SelectRow(ctx, conn, row, "WHERE principal_id = ?", pri.PrincipalID.String())
	} else {
		err = s.updatePrincipalRow(ctx, conn, pri, row)
	}
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.savePrincipalGroup(ctx, conn, pri, row)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.savePrincipalRole(ctx, conn, pri, row)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Service) DeletePrincipal(ctx Context, conn seacle.Executable, pri *entity.Principal) error {
	row := &PrincipalRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE principal_id = ?", pri.PrincipalID.String())
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	if err == sql.ErrNoRows {
		return nil
	} else {
		err = s.deletePrincipalRow(ctx, conn, pri, row)
	}

	// TODO delete PrincipalGroup / PrincipalRole

	return err
}

func (s *Service) createPrincipalRow(ctx Context, conn seacle.Executable, pri *entity.Principal) error {
	row := principalRowFromEntity(pri)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updatePrincipalRow(ctx Context, conn seacle.Executable, pri *entity.Principal, row *PrincipalRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
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

func (s *Service) deletePrincipalRow(ctx Context, conn seacle.Executable, pri *entity.Principal, row *PrincipalRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock principal row: err=%s", err)
	}

	// delete row
	err = seacle.Delete(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update principal row: err=%s", err)
	}

	return nil
}

func (s *Service) savePrincipalRole(ctx Context, conn seacle.Executable, pri *entity.Principal, priRow *PrincipalRow) error {
	roles := pri.AttachedRoles()
	roleSeqIDs := make([]int64, 0, len(roles))
	if len(roles) != 0 {
		roleIDs := make([]string, 0, len(roles))
		for _, v := range roles {
			roleIDs = append(roleIDs, v.RoleID.String())
		}
		roleRows := []*RoleRow{}
		err := seacle.Select(ctx, conn, &roleRows, `WHERE role_id IN (?)`, roleIDs)
		if err != nil {
			return err
		}
		for _, v := range roleRows {
			roleSeqIDs = append(roleSeqIDs, v.SeqID)
		}
	}

	priRoleMaps := []*PrincipalRoleMapRow{}
	err := seacle.Select(ctx, conn, &priRoleMaps, `WHERE principal_seq_id = ?`, priRow.SeqID)
	if err != nil {
		return err
	}
	args := make([]relationRow, 0, len(priRoleMaps))
	for _, v := range priRoleMaps {
		args = append(args, v)
	}

	added, deleted := compareSeqID(roleSeqIDs, args)
	if len(added) != 0 {
		for _, roleSeqID := range added {
			_, err = seacle.Insert(ctx, conn, &PrincipalRoleMapRow{
				PrincipalSeqID: priRow.SeqID,
				RoleSeqID:      roleSeqID,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	if len(deleted) != 0 {
		for _, roleSeqID := range deleted {
			err = seacle.Delete(ctx, conn, &PrincipalRoleMapRow{
				PrincipalSeqID: priRow.SeqID,
				RoleSeqID:      roleSeqID,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (s *Service) savePrincipalGroup(ctx Context, conn seacle.Executable, pri *entity.Principal, priRow *PrincipalRow) error {
	groups := pri.Groups()
	groupSeqIDs := make([]int64, 0, len(groups))
	if len(groups) != 0 {
		groupIDs := make([]string, 0, len(groups))
		for _, v := range groups {
			groupIDs = append(groupIDs, v.GroupID.String())
		}
		groupRows := []*GroupRow{}
		err := seacle.Select(ctx, conn, &groupRows, `WHERE group_id IN (?)`, groupIDs)
		if err != nil {
			return err
		}
		for _, v := range groupRows {
			groupSeqIDs = append(groupSeqIDs, v.SeqID)
		}
	}

	priGroupMaps := []*PrincipalGroupMapRow{}
	err := seacle.Select(ctx, conn, &priGroupMaps, `WHERE principal_seq_id = ?`, priRow.SeqID)
	if err != nil {
		return err
	}
	args := make([]relationRow, 0, len(priGroupMaps))
	for _, v := range priGroupMaps {
		args = append(args, v)
	}

	added, deleted := compareSeqID(groupSeqIDs, args)
	if len(added) != 0 {
		for _, groupSeqID := range added {
			_, err = seacle.Insert(ctx, conn, &PrincipalGroupMapRow{
				PrincipalSeqID: priRow.SeqID,
				GroupSeqID:     groupSeqID,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	if len(deleted) != 0 {
		for _, groupSeqID := range deleted {
			err = seacle.Delete(ctx, conn, &PrincipalGroupMapRow{
				PrincipalSeqID: priRow.SeqID,
				GroupSeqID:     groupSeqID,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (s *Service) FindPrincipals(ctx Context, conn seacle.Selectable, principalIDs []string) ([]*entity.Principal, error) {
	if len(principalIDs) == 0 {
		return []*entity.Principal{}, nil
	}

	pris := make([]*PrincipalRow, 0, 8)
	err := seacle.Select(ctx, conn, &pris, `WHERE principal_id IN (?) ORDER BY seq_id`, principalIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(pris) == 0 {
		return []*entity.Principal{}, nil
	}

	// seq_idを抽出
	priSeqIDs := make([]int64, 0, len(pris))
	for _, v := range pris {
		priSeqIDs = append(priSeqIDs, v.SeqID)
	}

	// AuthOIDC
	oidcs := make([]*AuthOIDCRow, 0, len(pris))
	authMap := map[int64]*entity.OIDCAuthorization{}
	err = seacle.Select(ctx, conn, &oidcs, `WHERE principal_seq_id IN (?)`, priSeqIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, v := range oidcs {
		authMap[v.PrincipalSeqID] = v.ToEntity()
	}

	// APIKey
	apikeys := []*AuthAPIKeyRow{}
	apikeyMap := map[int64][]*entity.AuthAPIKey{}
	err = seacle.Select(ctx, conn, &apikeys, `WHERE principal_seq_id IN (?)`, priSeqIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, v := range apikeys {
		if apikeyMap[v.PrincipalSeqID] == nil {
			apikeyMap[v.PrincipalSeqID] = []*entity.AuthAPIKey{}
		}
		apikeyMap[v.PrincipalSeqID] = append(apikeyMap[v.PrincipalSeqID], v.ToEntity())
	}

	// Role
	rolesMap, err := s.findRolesByPrincipalSeqID(ctx, conn, priSeqIDs)
	if err != nil {
		return nil, err
	}

	// Group
	groupsMap, err := s.findGroupsByPrincipalSeqID(ctx, conn, priSeqIDs)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Principal, 0, len(pris))
	f := entity.NewFactory(s.q)
	for _, v := range pris {
		result = append(result, f.NewPrincipal(
			uuid.MustParse(v.PrincipalID), v.Name, v.Description,
			authMap[v.SeqID], apikeyMap[v.SeqID], rolesMap[v.SeqID], groupsMap[v.SeqID]))
	}

	return result, nil
}

func (s *Service) FindPrincipalByOIDC(ctx Context, conn seacle.Selectable, issuer, subject string) (*entity.Principal, error) {
	pri := PrincipalRow{}
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
	err := seacle.Select(ctx, conn, &pris, "ORDER BY seq_id")
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
