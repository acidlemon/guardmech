// Code generated by seacle.Generator DO NOT EDIT
// About seacle: https://github.com/acidlemon/seacle
package db

import (
	"database/sql"

	"github.com/acidlemon/seacle"
)

var _ seacle.Mappable = (*AuthAPIKeyRow)(nil)

func (p *AuthAPIKeyRow) Table() string {
	return "auth_apikey"
}

func (p *AuthAPIKeyRow) Columns() []string {
	return []string{"auth_apikey.seq_id", "auth_apikey.auth_apikey_id", "auth_apikey.name", "auth_apikey.masked_token", "auth_apikey.hashed_token", "auth_apikey.principal_seq_id"}
}

func (p *AuthAPIKeyRow) PrimaryKeys() []string {
	return []string{"seq_id"}
}

func (p *AuthAPIKeyRow) PrimaryValues() []interface{} {
	return []interface{}{p.SeqID}
}

func (p *AuthAPIKeyRow) ValueColumns() []string {
	return []string{"auth_apikey_id", "name", "masked_token", "hashed_token", "principal_seq_id"}
}

func (p *AuthAPIKeyRow) Values() []interface{} {
	return []interface{}{p.AuthAPIKeyID, p.Name, p.MaskedToken, p.HashedToken, p.PrincipalSeqID}
}

func (p *AuthAPIKeyRow) AutoIncrementColumn() string {
	return "seq_id"
}

func (p *AuthAPIKeyRow) Scan(r seacle.RowScanner) error {
	var arg0 int64
	var arg1 string
	var arg2 string
	var arg3 string
	var arg4 string
	var arg5 int64

	err := r.Scan(&arg0, &arg1, &arg2, &arg3, &arg4, &arg5)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		// something wrong
		return err
	}

	p.SeqID = arg0
	p.AuthAPIKeyID = arg1
	p.Name = arg2
	p.MaskedToken = arg3
	p.HashedToken = arg4
	p.PrincipalSeqID = arg5

	return nil
}
