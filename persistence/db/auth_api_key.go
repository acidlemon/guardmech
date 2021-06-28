package db

type AuthAPIKeyRow struct {
	SeqID          int64  `db:"seq_id,primary"`
	AuthAPIKey     string `db:"auth_apikey_id"`
	Name           string `db:"name"`
	MaskedToken    string `db:"masked_token"`
	Salt           string `db:"salt"`
	HashedToken    string `db:"hashed_token"`
	PrincipalSeqID int64  `db:"principal_seq_id"`
}

/*
func (s *Service) SaveAPIKey(ctx Context, tx *db.Tx, a *membership.APIKey) (int64, error) {
	if a.SeqID == 0 {
		return s.createAPIKey(ctx, tx, a)
	} else {
		return s.updateAPIKey(ctx, tx, a)
	}
}

func (s *Service) createAPIKey(ctx Context, tx *db.Tx, a *membership.APIKey) (int64, error) {
	ap := AuthAPIKeyRow{
		*a, a.Principal.SeqID,
	}
	seqID, err := seacle.Insert(ctx, tx, &ap)
	if err != nil {
		return 0, err
	}
	return seqID, nil
}

func (s *Service) updateAPIKey(ctx Context, tx *db.Tx, a *membership.APIKey) (int64, error) {
	apForUpdate := &AuthAPIKeyRow{}
	err := seacle.SelectRow(ctx, tx, apForUpdate, `WHERE unique_id = ?`, a.UniqueID)
	if err != nil {
		// TODO: fallback to createAPIKey?
		return 0, nil
	}

	if apForUpdate.SeqID != a.SeqID {
		// ???
		return 0, fmt.Errorf("ID mismatched. seqID=%d, a.SeqID=%d", apForUpdate.SeqID, a.SeqID)
	}

	// lock row
	err = seacle.SelectRow(ctx, tx, apForUpdate, `FOR UPDATE WHERE seq_id = ?`, apForUpdate.SeqID)
	if err != nil {
		return 0, fmt.Errorf("failed to lock api_key row: err=%s", err)
	}

	// update row
	ap := AuthAPIKeyRow{
		*a, a.Principal.SeqID,
	}
	err = seacle.Update(ctx, tx, &ap)
	if err != nil {
		return 0, fmt.Errorf("failed to update api_key row: err=%s", err)
	}

	return a.SeqID, nil
}
*/
