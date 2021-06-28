// Code generated by seacle.Generator DO NOT EDIT
// About seacle: https://github.com/acidlemon/seacle
package db

import (
	"database/sql"

	"github.com/acidlemon/seacle"
)

var _ seacle.Mappable = (*GroupRow)(nil)

func (p *GroupRow) Table() string {
	return "group_info"
}

func (p *GroupRow) Columns() []string {
	return []string{"group_info.seq_id", "group_info.group_id", "group_info.name", "group_info.description"}
}

func (p *GroupRow) PrimaryKeys() []string {
	return []string{"seq_id"}
}

func (p *GroupRow) PrimaryValues() []interface{} {
	return []interface{}{p.SeqID}
}

func (p *GroupRow) ValueColumns() []string {
	return []string{"group_id", "name", "description"}
}

func (p *GroupRow) Values() []interface{} {
	return []interface{}{p.GroupID, p.Name, p.Description}
}

func (p *GroupRow) Scan(r seacle.RowScanner) error {
	var arg0 int64
	var arg1 string
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
	p.GroupID = arg1
	p.Name = arg2
	p.Description = arg3

	return nil
}
