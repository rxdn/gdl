package application

type InstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions uint64   `json:"permissions,string"`
}
