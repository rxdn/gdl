package user

type ClientStatus struct {
	Desktop ClientStatusType `json:"desktop"`
	Mobile  ClientStatusType `json:"mobile"`
	Web     ClientStatusType `json:"web"`
}
