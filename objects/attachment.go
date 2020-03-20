package objects

type Attachment struct {
	Id       uint64 `json:",string"`
	Filename string
	Size     int
	url      string
	ProxyUrl string
	height   int
	Width    int
}
