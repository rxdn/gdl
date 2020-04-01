package channel

type Attachment struct {
	Id       uint64 `json:",string"`
	Filename string
	Size     int
	Url      string
	ProxyUrl string
	height   int
	Width    int
}
