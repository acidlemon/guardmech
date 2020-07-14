package infra

/*
type RowScanner interface {
	Scan(dest ...interface{}) error
}

func scanPrincipalRow(r RowScanner) (*membership.Principal, error) {
	var id int64
	var uuid uuid.UUID
	var name, description string
	err := r.Scan(&id, &uuid, &name, &description)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &membership.Principal{
		SeqID:       id,
		UniqueID:    uuid,
		Name:        name,
		Description: description,
	}, nil
}

func scanAuthRow(r RowScanner) (*membership.Auth, error) {
	var seqID int64
	var uniqID uuid.UUID
	var issuer, subject, email string
	err := r.Scan(&seqID, &uniqID, &issuer, &subject, &email)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &membership.Auth{
		SeqID:    seqID,
		UniqueID: uniqID,
		Issuer:   issuer,
		Subject:  subject,
		Email:    email,
	}, nil
}

func scanAPIKeyRow(r RowScanner) (*membership.APIKey, error) {
	var seqID int64
	var uniqID uuid.UUID
	var token string
	err := r.Scan(&seqID, &uniqID, &token)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		// something wrong
		return nil, err
	}

	return &membership.APIKey{
		SeqID:    seqID,
		UniqueID: uniqID,
		Token:    token,
	}, nil
}

func scanRoleRow(r RowScanner) (*membership.Role, error) {
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

	return &membership.Role{
		SeqID:       seqID,
		UniqueID:    uniqID,
		Name:        name,
		Description: description,
	}, nil
}

func scanPermissionRow(r RowScanner) (*membership.Permission, error) {
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

	return &membership.Permission{
		SeqID:       seqID,
		UniqueID:    uniqID,
		Name:        name,
		Description: description,
	}, nil
}
*/
