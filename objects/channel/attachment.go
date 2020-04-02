package channel

type Attachment struct {
	Id       uint64 `json:",string"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}
