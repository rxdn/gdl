package objects

type Embed struct {
	Title       string
	Type        string
	Description string
	Url         string
	Timestamp   string
	Color       int
	Footer      EmbedField
	Image       EmbedImage
	Thumbnail   EmbedThumbnail
	Video       EmbedVideo
	Provider    EmbedProvider
	Author      EmbedAuthor
	Fields      []EmbedField
}

type EmbedAuthor struct {
	Name         string
	Url          string
	IconUrl      string
	ProxyIconUrl string
}

type EmbedField struct {
	Name   string
	Value  string
	Inline bool
}

type EmbedFooter struct {
	Text         string
	IconUrl      string
	ProxyIconUrl string
}

type EmbedImage struct {
	Url      string
	ProxyUrl string
	Height   int
	Width    int
}

type EmbedProvider struct {
	Name string
	Url  string
}

type EmbedThumbnail struct {
	Url      string
	ProxyUrl string
	Height   int
	Width    int
}

type EmbedVideo struct {
	Url      string
	Height   int
	Width    int
}
