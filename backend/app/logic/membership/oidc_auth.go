package membership

import "github.com/google/uuid"

type OIDCAuthorization struct {
	OIDCAuthID uuid.UUID
	Issuer     string
	Subject    string
	Email      string
	Name       string
}
