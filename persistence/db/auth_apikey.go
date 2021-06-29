package db

import (
	"database/sql"
	"fmt"
	"log"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type AuthAPIKeyRow struct {
	SeqID          int64  `db:"seq_id,primary,auto_increment"`
	AuthAPIKeyID   string `db:"auth_apikey_id"`
	Name           string `db:"name"`
	MaskedToken    string `db:"masked_token"`
	HashedToken    string `db:"hashed_token"`
	PrincipalSeqID int64  `db:"principal_seq_id"`
}

func authAPIKeyRowFromEntity(a *entity.AuthAPIKey, principalSeqID int64) *AuthAPIKeyRow {
	return &AuthAPIKeyRow{
		AuthAPIKeyID:   a.AuthAPIKeyID.String(),
		Name:           a.Name,
		MaskedToken:    a.MaskedToken,
		HashedToken:    a.HashedToken,
		PrincipalSeqID: principalSeqID,
	}
}

func (a *AuthAPIKeyRow) ToEntity() *entity.AuthAPIKey {
	return &entity.AuthAPIKey{
		AuthAPIKeyID: uuid.MustParse(a.AuthAPIKeyID),
		Name:         a.Name,
		MaskedToken:  a.MaskedToken,
		HashedToken:  a.HashedToken,
	}
}

func (s *Service) SaveAuthAPIKey(ctx Context, conn seacle.Executable, a *entity.AuthAPIKey, pri *entity.Principal) error {
	priRow := &PrincipalRow{}
	err := seacle.SelectRow(ctx, conn, priRow, "WHERE principal_id = ?", pri.PrincipalID.String())
	if err != nil && err != sql.ErrNoRows {
		log.Println("failed to select parent principal")
		return err
	}

	row := &AuthAPIKeyRow{}
	err = seacle.SelectRow(ctx, conn, row, "WHERE auth_apikey_id = ?", a.AuthAPIKeyID.String())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return s.createAuthAPIKey(ctx, conn, a, priRow.SeqID)
	} else {
		return s.updateAuthAPIKey(ctx, conn, a, priRow.SeqID, row)
	}
}

func (s *Service) createAuthAPIKey(ctx Context, conn seacle.Executable, a *entity.AuthAPIKey, priSeqID int64) error {
	row := authAPIKeyRowFromEntity(a, priSeqID)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) updateAuthAPIKey(ctx Context, conn seacle.Executable, a *entity.AuthAPIKey, priSeqID int64, row *AuthAPIKeyRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock api_key row: err=%s", err)
	}

	// update row
	row.Name = a.Name
	row.MaskedToken = a.MaskedToken
	row.HashedToken = a.HashedToken
	row.PrincipalSeqID = priSeqID
	err = seacle.Update(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update api_key row: err=%s", err)
	}

	return nil
}
