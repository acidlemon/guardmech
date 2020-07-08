package infra

import "github.com/acidlemon/guardmech/membership"

type APIKey membership.APIKey

/*
func (s *Membership) SaveAPIKey(ctx context.Context, tx *db.Tx, key *membership.APIKey) (int64, error) {
	if key.SeqID == 0 {
		return s.createAPIKey(ctx, tx, key)
	} else {
		return s.updateAPIKey(ctx, tx, key)
	}
}
*/
/*
   seq_id BIGINT NOT NULL auto_increment,
   unique_id VARCHAR(40) CHARACTER SET latin1 NOT NULL UNIQUE,
   principal_id BIGINT NOT NULL,
   name VARCHAR(191) NOT NULL,
   masked_token VARCHAR(255) CHARACTER SET utf8 NOT NULL,
   salt VARCHAR(255) CHARACTER SET utf8 NOT NULL,
   hashed_token VARCHAR(255) CHARACTER SET utf8 NOT NULL UNIQUE,
*/
/*
func (s *Membership) createAPIKey(ctx context.Context, tx *db.Tx, key *membership.APIKey) (int64, error) {
	// new Principal
	res, err := tx.ExecContext(ctx,
		`INSERT INTO api_key (unique_id, name, description) VALUES (?, ?, ?)`,
		pri.UniqueID, pri.Name, pri.Description,
	)
	if err != nil {
		return 0, err
	}
	seqID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return seqID, nil
}

func (s *Membership) updateAPIKey(ctx context.Context, tx *db.Tx, key *membership.APIKey) (int64, error) {

}
*/
