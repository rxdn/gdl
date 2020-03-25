package objects

import "github.com/rxdn/gdl/permission"

type PermissionOverwriteType string

const (
	PermissionTypeRole   PermissionOverwriteType = "role"
	PermissionTypeMember PermissionOverwriteType = "member"
)

type PermissionOverwrite struct {
	Id    uint64                  `json:"id,string,omitempty"`
	Type  PermissionOverwriteType `json:"type"`
	Allow int                     `json:"allow"`
	Deny  int                     `json:"deny"`
}

func (o *PermissionOverwrite) GetAllowedPermissions() []permission.Permission {
	perms := make([]permission.Permission, 0)

	for _, perm := range permission.AllPermissions {
		if permission.HasPermission(o.Allow, perm) {
			perms = append(perms, perm)
		}
	}

	return perms
}

func (o *PermissionOverwrite) GetDeniedPermissions() []permission.Permission {
	perms := make([]permission.Permission, 0)

	for _, perm := range permission.AllPermissions {
		if permission.HasPermission(o.Deny, perm) {
			perms = append(perms, perm)
		}
	}

	return perms
}
