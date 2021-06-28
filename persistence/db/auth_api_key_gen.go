// Code generated by seacle.Generator DO NOT EDIT
// About seacle: https://github.com/acidlemon/seacle
package db

import (
	"database/sql"

	"github.com/acidlemon/seacle"
)

var _ seacle.Mappable = (*AuthAPIKeyRow)(nil)

func (p *AuthAPIKeyRow) Table() string {
	return "auth_api_key"
}

func (p *AuthAPIKeyRow) Columns() []string {
	return []string{"auth_api_key.seq_id", "auth_api_key.auth_apikey_id", "auth_api_key.name", "auth_api_key.masked_token", "auth_api_key.salt", "auth_api_key.hashed_token", "auth_api_key.principal_seq_id"}
}

func (p *AuthAPIKeyRow) PrimaryKeys() []string {
	return []string{"seq_id"}
}

func (p *AuthAPIKeyRow) PrimaryValues() []interface{} {
	return []interface{}{p.SeqID}
}

func (p *AuthAPIKeyRow) ValueColumns() []string {
	return []string{"auth_apikey_id", "name", "masked_token", "salt", "hashed_token", "principal_seq_id"}
}

func (p *AuthAPIKeyRow) Values() []interface{} {
	return []interface{}{p.AuthAPIKey, p.Name, p.MaskedToken, p.Salt, p.HashedToken, p.PrincipalSeqID}
}

func (p *AuthAPIKeyRow) Scan(r seacle.RowScanner) error {
	var arg0 int64
	var arg1 string
	var arg2 string
	var arg3 string
	var arg4 string
	var arg5 string
	var arg6 int64

	err := r.Scan(&arg0, &arg1, &arg2, &arg3, &arg4, &arg5, &arg6)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		// something wrong
		return err
	}

	p.SeqID = arg0
	p.AuthAPIKey = arg1
	p.Name = arg2
	p.MaskedToken = arg3
	p.Salt = arg4
	p.HashedToken = arg5
	p.PrincipalSeqID = arg6

	return nil
}
