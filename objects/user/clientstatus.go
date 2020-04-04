package user

type ClientStatus struct {
	Desktop ClientStatusType `json:"desktop"`
	Mobile  string           `json:"mobile"`
	Web     string           `json:"web"`
}
