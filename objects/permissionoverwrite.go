package objects

type PermissionOverwrite struct {
	Id    uint64                  `json:"id,string"`
	Type  PermissionOverwriteType `json:"type"`
	Allow int                     `json:"allow"`
	Deny  int                     `json:"deny"`
}
