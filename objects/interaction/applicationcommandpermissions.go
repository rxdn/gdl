package interaction

type ApplicationCommandPermissions struct {
	Id         uint64                           `json:"id,string"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}
