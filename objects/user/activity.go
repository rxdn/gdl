package user

type Activity struct {
	Name          string       `json:"name"`
	Type          ActivityType `json:"type"`
	Url           string       `json:"url"`
	Timestamps    Timestamps   `json:"timestamps"`
	ApplicationId uint64       `json:"application_id,string"`
	Details       string       `json:"details"`
	State         string       `json:"state"`
	// TODO: Figure out how to handle emoji w/o import cycle
	Party    Party  `json:"party"`
	Assets   Asset  `json:"assets"`
	Secret   Secret `json:"secret"`
	Instance bool   `json:"instance"`
	Flags    int    `json:"flags"` // TODO: Wrap this
}
