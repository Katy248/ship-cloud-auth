package data

import (
	"slices"

	"github.com/google/uuid"
)

const (
	PermissionUserEdit    = "user.edit"
	PermissionUserList    = "user.list"
	PermissionUserBlock   = "user.block"
	PermissionUserGetByID = "user.get-by-id"
)

func (s Session) HasPermission(p string) bool {
	return slices.Contains(s.Permissions, p)
}

func (s Session) CanGetUserByID(userID uuid.UUID) bool {
	if s.UserID == userID {
		return true
	}
	return s.HasPermission(PermissionUserGetByID)
}

func (s Session) CanEditUser(userID uuid.UUID) bool {
	if s.UserID == userID {
		return true
	}
	return s.HasPermission(PermissionUserEdit)
}
