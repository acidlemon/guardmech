package membership

import "github.com/google/uuid"

type APIKey struct {
	SeqID       int64      `json:"seq_id"`
	UniqueID    uuid.UUID  `json:"unique_id"`
	MaskedToken string     `json:"token" db:"masked_token"`
	Salt        string     `json:"-" db:"salt"`
	HashedToken string     `json:"-" db:"hashed_token"`
	Principal   *Principal `json:"-" db:"-"`
}

/*
   masked_token VARCHAR(255) CHARACTER SET utf8 NOT NULL,
   salt VARCHAR(255) CHARACTER SET utf8 NOT NULL,
   hashed_token VARCHAR(255) CHARACTER SET utf8 NOT NULL UNIQUE,
*/
