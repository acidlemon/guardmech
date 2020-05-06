package infra

import (
	"database/sql"

	"github.com/acidlemon/guardmech/membership"
	"github.com/google/uuid"
)

func scanGroupRow(r RowScanner) (*membership.Group, error) {
	var seqID int64
	var uniqID uuid.UUID
	var name, description string
	err := r.Scan(&seqID, &uniqID, &name, &description)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &membership.Group{
		SeqID:       seqID,
		UniqueID:    uniqID,
		Name:        name,
		Description: description,
	}, nil
}
