// Code generated by seacle.Generator DO NOT EDIT
// About seacle: https://github.com/acidlemon/seacle
package db

import (
	"database/sql"

	"github.com/acidlemon/seacle"
)

var _ seacle.Mappable = (*PrincipalGroupMapRow)(nil)

func (p *PrincipalGroupMapRow) Table() string {
	return "principal_group_map"
}

func (p *PrincipalGroupMapRow) Columns() []string {
	return []string{"principal_group_map.principal_seq_id", "principal_group_map.group_seq_id"}
}

func (p *PrincipalGroupMapRow) PrimaryKeys() []string {
	return []string{"principal_seq_id", "group_seq_id"}
}

func (p *PrincipalGroupMapRow) PrimaryValues() []interface{} {
	return []interface{}{p.PrincipalSeqID, p.GroupSeqID}
}

func (p *PrincipalGroupMapRow) ValueColumns() []string {
	return []string{}
}

func (p *PrincipalGroupMapRow) Values() []interface{} {
	return []interface{}{}
}

func (p *PrincipalGroupMapRow) AutoIncrementColumn() string {
	return ""
}

func (p *PrincipalGroupMapRow) Scan(r seacle.RowScanner) error {
	var arg0 int64
	var arg1 int64

	err := r.Scan(&arg0, &arg1)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		// something wrong
		return err
	}

	p.PrincipalSeqID = arg0
	p.GroupSeqID = arg1

	return nil
}
