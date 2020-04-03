package user

type UpdateStatus struct {
	Since  int      `json:"since,omitempty"` // time since client went idle (unix epoch millis)
	Game   Activity `json:"game,omitempty"`
	Status string   `json:"status"`
	Afk    bool     `json:"afk"`
}

func BuildStatus(activityType ActivityType, status string) UpdateStatus {
	return UpdateStatus{
		Game: Activity{
			Type: activityType,
		},
		Status: status,
	}
}
