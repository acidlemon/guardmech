// Code generated by seacle.Generator DO NOT EDIT
// About seacle: https://github.com/acidlemon/seacle
package infra

import (
	"database/sql"

	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

var _ seacle.Mappable = (*Role)(nil)

func (p *Role) Table() string {
	return "role_info AS role_info"
}

func (p *Role) Columns() []string {
	return []string{"role_info.seq_id", "role_info.unique_id", "role_info.name", "role_info.description"}
}

func (p *Role) PrimaryKeys() []string {
	return []string{"role_info.seq_id"}
}

func (p *Role) PrimaryValues() []interface{} {
	return []interface{}{p.SeqID}
}

func (p *Role) ValueColumns() []string {
	return []string{"role_info.unique_id", "role_info.name", "role_info.description"}
}

func (p *Role) Values() []interface{} {
	return []interface{}{p.UniqueID, p.Name, p.Description}
}

func (p *Role) Scan(r seacle.RowScanner) error {
	var arg0 int64
	var arg1 uuid.UUID
	var arg2 string
	var arg3 string

	err := r.Scan(&arg0, &arg1, &arg2, &arg3)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		// something wrong
		return err
	}

	p.SeqID = arg0
	p.UniqueID = arg1
	p.Name = arg2
	p.Description = arg3

	return nil
}
