package interaction

type ApplicationCommandPermissionType uint8

const (
	ApplicationCommandPermissionTypeRole ApplicationCommandPermissionType = iota + 1
	ApplicationCommandPermissionTypeUser
)
