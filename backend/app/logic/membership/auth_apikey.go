package membership

import "github.com/google/uuid"

type AuthAPIKey struct {
	AuthAPIKeyID uuid.UUID
	Name         string
	MaskedToken  string
	HashedToken  string
}
