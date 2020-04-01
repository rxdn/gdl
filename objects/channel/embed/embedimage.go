package embed

type EmbedImage struct {
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}
