package objects

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
