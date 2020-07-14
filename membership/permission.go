package membership

import (
	"github.com/google/uuid"
)

const (
	PermissionOwner = "_GUARDMECH_OWNER"
)

type Permission struct {
	SeqID       int64
	UniqueID    uuid.UUID `json:"unique_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
