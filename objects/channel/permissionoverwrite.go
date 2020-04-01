package channel

type PermissionOverwrite struct {
	Id    uint64                  `json:"id,string,omitempty"`
	Type  PermissionOverwriteType `json:"type"`
	Allow int                     `json:"allow"`
	Deny  int                     `json:"deny"`
}
