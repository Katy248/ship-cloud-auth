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

	PermissionOrganizationList             = "organization.list"
	PermissionOrganizationEdit             = "organization.edit"
	PermissionOrganizationGetByID          = "organization.get-by-id"
	PermissionOrganizationDelete           = "organization.delete"
	PermissionOrganizationInviteMemeber    = "organization.members.invite"
	PermissionOrganizationListMemebers     = "organization.members.list"
	PermissionOrganizationRemoveMemebers   = "organization.members.remove"
	PermissionOrganizationCreateRole       = "organization.roles.create"
	PermissionOrganizationEditRoles        = "organization.roles.edit"
	PermissionOrganizationDeleteRoles      = "organization.roles.delete"
	PermissionOrganizationConnectDevice    = "organization.devices.connect"
	PermissionOrganizationDisconnectDevice = "organization.devices.disconnect"
	PermissionOrganizationDeviceEditTags   = "organization.devices.edit-tags"
	PermissionOrganizationTagsEdit         = "organization.tags.edit"
)

func GetAllPermissions() []string {
	return []string{
		PermissionUserEdit,
		PermissionUserList,
		PermissionUserBlock,
		PermissionUserGetByID,

		PermissionOrganizationList,
		PermissionOrganizationEdit,
		PermissionOrganizationGetByID,
		PermissionOrganizationDelete,
		PermissionOrganizationInviteMemeber,
		PermissionOrganizationListMemebers,
		PermissionOrganizationRemoveMemebers,
		PermissionOrganizationCreateRole,
		PermissionOrganizationEditRoles,
		PermissionOrganizationDeleteRoles,
		PermissionOrganizationDeviceEditTags,
		PermissionOrganizationConnectDevice,
		PermissionOrganizationDisconnectDevice,
		PermissionOrganizationTagsEdit,
	}
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
