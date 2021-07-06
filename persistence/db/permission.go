package db

import (
	"database/sql"
	"fmt"
	"log"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type PermissionRow struct {
	SeqID        int64  `db:"seq_id,primary,auto_increment"`
	PermissionID string `db:"permission_id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
}

func permissionRowFromEntity(perm *entity.Permission) *PermissionRow {
	return &PermissionRow{
		PermissionID: perm.PermissionID.String(),
		Name:         perm.Name,
		Description:  perm.Description,
	}
}

func (perm *PermissionRow) ToEntity() *entity.Permission {
	return &entity.Permission{
		PermissionID: uuid.MustParse(perm.PermissionID),
		Name:         perm.Name,
		Description:  perm.Description,
	}
}

func (s *Service) SavePermission(ctx Context, conn seacle.Executable, perm *entity.Permission) error {
	row := &PermissionRow{}
	err := seacle.SelectRow(ctx, conn, row, "WHERE permission_id = ?", perm.PermissionID.String())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return s.createPermission(ctx, conn, perm)
	} else {
		return s.updatePermission(ctx, conn, perm, row)
	}
}

func (s *Service) createPermission(ctx Context, conn seacle.Executable, perm *entity.Permission) error {
	row := permissionRowFromEntity(perm)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updatePermission(ctx Context, conn seacle.Executable, perm *entity.Permission, row *PermissionRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock permission row: err=%s", err)
	}

	// update row
	row.Name = perm.Name
	row.Description = perm.Description
	err = seacle.Update(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update permission row: err=%s", err)
	}

	return nil
}

func (s *Service) FindPermissions(ctx Context, conn seacle.Selectable, permissionIDs []string) ([]*entity.Permission, error) {
	if len(permissionIDs) == 0 {
		return []*entity.Permission{}, nil
	}

	permRows := make([]*PermissionRow, 0, 8)
	err := seacle.Select(ctx, conn, &permRows, `WHERE permission_id IN (?)`, permissionIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	perms := []*entity.Permission{}
	for _, v := range permRows {
		perms = append(perms, v.ToEntity())
	}
	return perms, nil
}

func (s *Service) EnumeratePermissionIDs(ctx Context, conn seacle.Selectable) ([]string, error) {
	perms := make([]*PermissionRow, 0, 8)
	err := seacle.Select(ctx, conn, &perms, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]string, 0, len(perms))
	for _, v := range perms {
		result = append(result, v.PermissionID)
	}

	return result, nil
}
