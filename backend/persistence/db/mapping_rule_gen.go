// Code generated by seacle.Generator DO NOT EDIT
// About seacle: https://github.com/acidlemon/seacle
package db

import (
	"database/sql"

	"github.com/acidlemon/seacle"
)

var _ seacle.Mappable = (*MappingRuleRow)(nil)

func (p *MappingRuleRow) Table() string {
	return "mapping_rule"
}

func (p *MappingRuleRow) Columns() []string {
	return []string{"mapping_rule.seq_id", "mapping_rule.mapping_rule_id", "mapping_rule.rule_type", "mapping_rule.detail", "mapping_rule.name", "mapping_rule.description", "mapping_rule.priority", "mapping_rule.association_type", "mapping_rule.association_id"}
}

func (p *MappingRuleRow) PrimaryKeys() []string {
	return []string{"seq_id"}
}

func (p *MappingRuleRow) PrimaryValues() []interface{} {
	return []interface{}{p.SeqID}
}

func (p *MappingRuleRow) ValueColumns() []string {
	return []string{"mapping_rule_id", "rule_type", "detail", "name", "description", "priority", "association_type", "association_id"}
}

func (p *MappingRuleRow) Values() []interface{} {
	return []interface{}{p.MappingRuleID, p.RuleType, p.Detail, p.Name, p.Description, p.Priority, p.AssociationType, p.AssociationID}
}

func (p *MappingRuleRow) AutoIncrementColumn() string {
	return "seq_id"
}

func (p *MappingRuleRow) Scan(r seacle.RowScanner) error {
	var arg0 int64
	var arg1 string
	var arg2 int
	var arg3 string
	var arg4 string
	var arg5 string
	var arg6 int
	var arg7 int
	var arg8 string

	err := r.Scan(&arg0, &arg1, &arg2, &arg3, &arg4, &arg5, &arg6, &arg7, &arg8)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		// something wrong
		return err
	}

	p.SeqID = arg0
	p.MappingRuleID = arg1
	p.RuleType = arg2
	p.Detail = arg3
	p.Name = arg4
	p.Description = arg5
	p.Priority = arg6
	p.AssociationType = arg7
	p.AssociationID = arg8

	return nil
}
