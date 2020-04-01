package embed

type EmbedThumbnail struct {
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}
