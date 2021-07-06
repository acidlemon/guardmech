package db

import (
	"database/sql"
	"fmt"
	"log"

	entity "github.com/acidlemon/guardmech/app/logic/membership"
	"github.com/acidlemon/seacle"
	"github.com/google/uuid"
)

type AuthOIDCRow struct {
	SeqID          int64  `db:"seq_id,primary,auto_increment"`
	AuthOIDCID     string `db:"auth_oidc_id"`
	Issuer         string `db:"issuer"`
	Subject        string `db:"subject"`
	Email          string `db:"email"`
	Name           string `db:"name"`
	PrincipalSeqID int64  `db:"principal_seq_id"`
}

func authOIDCRowFromEntity(a *entity.OIDCAuthorization, principalSeqID int64) *AuthOIDCRow {
	return &AuthOIDCRow{
		AuthOIDCID:     a.OIDCAuthID.String(),
		Issuer:         a.Issuer,
		Subject:        a.Subject,
		Email:          a.Email,
		Name:           a.Name,
		PrincipalSeqID: principalSeqID,
	}
}

func (a *AuthOIDCRow) ToEntity() *entity.OIDCAuthorization {
	return &entity.OIDCAuthorization{
		OIDCAuthID: uuid.MustParse(a.AuthOIDCID),
		Issuer:     a.Issuer,
		Subject:    a.Subject,
		Email:      a.Email,
		Name:       a.Name,
	}
}

func (s *Service) SaveAuthOIDC(ctx Context, conn seacle.Executable, a *entity.OIDCAuthorization, pri *entity.Principal) error {
	priRow := &PrincipalRow{}
	err := seacle.SelectRow(ctx, conn, priRow, "WHERE principal_id = ?", pri.PrincipalID.String())
	if err != nil && err != sql.ErrNoRows {
		log.Println("failed to select parent principal")
		return err
	}

	row := &AuthOIDCRow{}
	err = seacle.SelectRow(ctx, conn, row, "WHERE auth_oidc_id = ?", a.OIDCAuthID.String())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return s.createAuthOIDC(ctx, conn, a, priRow.SeqID)
	} else {
		return s.updateAuthOIDC(ctx, conn, a, priRow.SeqID, row)
	}
}

func (s *Service) createAuthOIDC(ctx Context, conn seacle.Executable, a *entity.OIDCAuthorization, priSeqID int64) error {
	row := authOIDCRowFromEntity(a, priSeqID)
	_, err := seacle.Insert(ctx, conn, row)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) updateAuthOIDC(ctx Context, conn seacle.Executable, a *entity.OIDCAuthorization, priSeqID int64, row *AuthOIDCRow) error {
	// lock row
	err := seacle.SelectRow(ctx, conn, row, `WHERE seq_id = ? FOR UPDATE`, row.SeqID)
	if err != nil {
		return fmt.Errorf("failed to lock auth row: err=%s", err)
	}

	// update row
	row.Issuer = a.Issuer
	row.Subject = a.Subject
	row.Email = a.Email
	row.Name = a.Name
	row.PrincipalSeqID = priSeqID
	err = seacle.Update(ctx, conn, row)
	if err != nil {
		return fmt.Errorf("failed to update auth row: err=%s", err)
	}

	return nil
}
