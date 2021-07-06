// Code generated by seacle.Generator DO NOT EDIT
// About seacle: https://github.com/acidlemon/seacle
package db

import (
	"database/sql"

	"github.com/acidlemon/seacle"
)

var _ seacle.Mappable = (*RolePermissionMapRow)(nil)

func (p *RolePermissionMapRow) Table() string {
	return "role_permission_map"
}

func (p *RolePermissionMapRow) Columns() []string {
	return []string{"role_permission_map.role_seq_id", "role_permission_map.permission_seq_id"}
}

func (p *RolePermissionMapRow) PrimaryKeys() []string {
	return []string{"role_seq_id", "permission_seq_id"}
}

func (p *RolePermissionMapRow) PrimaryValues() []interface{} {
	return []interface{}{p.RoleSeqID, p.PermissionSeqID}
}

func (p *RolePermissionMapRow) ValueColumns() []string {
	return []string{}
}

func (p *RolePermissionMapRow) Values() []interface{} {
	return []interface{}{}
}

func (p *RolePermissionMapRow) AutoIncrementColumn() string {
	return ""
}

func (p *RolePermissionMapRow) Scan(r seacle.RowScanner) error {
	var arg0 int64
	var arg1 int64

	err := r.Scan(&arg0, &arg1)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		// something wrong
		return err
	}

	p.RoleSeqID = arg0
	p.PermissionSeqID = arg1

	return nil
}
