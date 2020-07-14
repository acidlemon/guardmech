// Code generated by seacle.Generator DO NOT EDIT
// About seacle: https://github.com/acidlemon/seacle
package infra

import (
	"database/sql"

	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

var _ seacle.Mappable = (*Group)(nil)

func (p *Group) Table() string {
	return "group_info AS group_info"
}

func (p *Group) Columns() []string {
	return []string{"group_info.seq_id", "group_info.unique_id", "group_info.name", "group_info.description"}
}

func (p *Group) PrimaryKeys() []string {
	return []string{"group_info.seq_id"}
}

func (p *Group) PrimaryValues() []interface{} {
	return []interface{}{p.SeqID}
}

func (p *Group) ValueColumns() []string {
	return []string{"group_info.unique_id", "group_info.name", "group_info.description"}
}

func (p *Group) Values() []interface{} {
	return []interface{}{p.UniqueID, p.Name, p.Description}
}

func (p *Group) Scan(r seacle.RowScanner) error {
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
