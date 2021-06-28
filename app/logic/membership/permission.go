package membership

import (
	"github.com/google/uuid"
)

const (
	PermissionOwner   = "_GUARDMECH_OWNER"
	PermissionOwnerID = "d4b6dc0b-f282-4e9c-b8d7-518f61737f21"
)

type Permission struct {
	PermissionID uuid.UUID `json:"unique_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
}

func newPermission(name, description string) *Permission {
	permissionID := uuid.New()
	if name == PermissionOwner {
		permissionID = uuid.MustParse(PermissionOwnerID)
	}

	return &Permission{
		PermissionID: permissionID,
		Name:         name,
		Description:  description,
	}
}
