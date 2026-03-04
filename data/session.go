package data

import (
	"slices"

	"github.com/google/uuid"
)

type Session struct {
	UserID      uuid.UUID `json:"userId"`
	Blocked     bool      `json:"blocked"`
	Permissions []string  `json:"permissions"`
}

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

const (
	PermissionUserEdit    = "user.edit"
	PermissionUserList    = "user.list"
	PermissionUserBlock   = "user.block"
	PermissionUserGetByID = "user.get-by-id"
)
