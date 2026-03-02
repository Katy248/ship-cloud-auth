package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPermissions(t *testing.T) {
	user := &User{
		Roles: []Role{
			{
				Permissions: fmt.Sprintf("%s:%s", BlockUserPermission, ListUsersPermission),
			},
			{
				Permissions: fmt.Sprintf("%s", CreateUserPermission),
			},
		},
	}
	assert.Equal(t, true, user.HasPermission(BlockUserPermission))
	assert.Equal(t, true, user.HasPermission(ListUsersPermission))
	assert.Equal(t, false, user.HasPermission(UnblockUserPermission))
	assert.Equal(t, true, user.HasPermission(CreateUserPermission))
}

func TestHasPermissions_all(t *testing.T) {
	user := &User{
		Roles: []Role{
			{
				Permissions: "all",
			},
		},
	}
	assert.Equal(t, true, user.HasPermission(BlockUserPermission))
	assert.Equal(t, true, user.HasPermission(ListUsersPermission))
	assert.Equal(t, true, user.HasPermission(UnblockUserPermission))
	assert.Equal(t, true, user.HasPermission(CreateUserPermission))
	assert.Equal(t, true, user.HasPermission(""))
	assert.Equal(t, true, user.HasPermission("sdsdsaddgsdrfoujgvfjklsdn"))
	assert.Equal(t, true, user.HasPermission("all"))
	assert.Equal(t, true, user.HasPermission("cool"))

}
