package membership

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

// APIKeyPrefix is the prefix of every raw API key token issued by guardmech.
const APIKeyPrefix = "gmch-"

type AuthAPIKey struct {
	AuthAPIKeyID uuid.UUID
	Name         string
	MaskedToken  string
	HashedToken  string
}

// HashAPIKeyToken returns the hex-encoded SHA-256 digest of a raw API key token.
// The digest is deterministic so that the hashed_token unique index can be used for lookup.
func HashAPIKeyToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
