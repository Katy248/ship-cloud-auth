package models

import (
	"slices"
	"strings"
)

const (
	AllPermissions = "all"

	ListUsersPermission   = "users.list"
	BlockUserPermission   = "users.block"
	UnblockUserPermission = "users.unblock"
	CreateUserPermission  = "users.create"
	UpdateUserPermission  = "users.update" // Deprecated: Should change to multimple permissions
	SetEmailPermission    = "users.set-email"
	SetPasswordPermission = "users.set-password"
)

func (u *User) HasPermission(permission string) bool {
	for _, role := range u.Roles {
		permissions := strings.Split(role.Permissions, ":")

		if slices.Contains(permissions, AllPermissions) {
			return true
		}

		if slices.Contains(permissions, permission) {
			return true
		}

	}
	return false
}
