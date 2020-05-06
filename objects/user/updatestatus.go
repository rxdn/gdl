package user

type UpdateStatus struct {
	Since  int              `json:"since"` // time since client went idle (unix epoch millis)
	Game   Activity         `json:"game,omitempty"`
	Status ClientStatusType `json:"status"`
	Afk    bool             `json:"afk"`
}

func BuildStatus(activityType ActivityType, status string) UpdateStatus {
	return UpdateStatus{
		Game: Activity{
			Name: status,
			Type: activityType,
		},
		Status: ClientStatusTypeOnline,
	}
}
