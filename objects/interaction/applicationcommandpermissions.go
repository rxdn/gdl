package interaction

type ApplicationCommandPermissionType uint8

const (
	ApplicationCommandPermissionTypeRole ApplicationCommandPermissionType = iota + 1
	ApplicationCommandPermissionTypeUser
)

type ApplicationCommandPermissions struct {
	Id         uint64                           `json:"id,string"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}
