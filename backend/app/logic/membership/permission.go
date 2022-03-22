package membership

import (
	"github.com/google/uuid"
)

const (
	PermissionOwnerName           = "_GUARDMECH_OWNER"
	PermissionOwnerDescription    = "Owner permission of guardmech"
	PermissionOwnerID             = "d4b6dc0b-f282-4e9c-b8d7-518f61737f21"
	PermissionWriteName           = "_GUARDMECH_WRITE"
	PermissionWriteDescription    = "Write permission of guardmech"
	PermissionWriteID             = "d4b6dc0b-f282-4e9c-b8d7-518f61737f22"
	PermissionReadOnlyName        = "_GUARDMECH_READONLY"
	PermissionReadOnlyDescription = "ReadOnly permission of guardmech"
	PermissionReadOnlyID          = "d4b6dc0b-f282-4e9c-b8d7-518f61737f23"
)

type Permission struct {
	PermissionID uuid.UUID
	Name         string
	Description  string
}

func newPermission(name, description string) *Permission {
	var permissionID uuid.UUID
	switch name {
	case PermissionOwnerName:
		permissionID = uuid.MustParse(PermissionOwnerID)

	case PermissionWriteName:
		permissionID = uuid.MustParse(PermissionWriteID)

	case PermissionReadOnlyName:
		permissionID = uuid.MustParse(PermissionReadOnlyID)

	default:
		permissionID = uuid.New()
	}

	return &Permission{
		PermissionID: permissionID,
		Name:         name,
		Description:  description,
	}
}
