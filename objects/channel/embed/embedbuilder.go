package embed

import "time"

func NewEmbed() *Embed {
	return new(Embed)
}

func (e *Embed) SetTitle(title string) *Embed {
	e.Title = title
	return e
}

func (e *Embed) SetDescription(description string) *Embed {
	e.Description = description
	return e
}

func (e *Embed) SetUrl(url string) *Embed {
	e.Url = url
	return e
}

func (e *Embed) SetTimestamp(timestamp time.Time) *Embed {
	e.Timestamp = &timestamp
	return e
}

func (e *Embed) SetColor(color int) *Embed {
	e.Color = color
	return e
}

func (e *Embed) SetFooter(text, iconUrl string) *Embed {
	e.Footer = EmbedFooter{
		Text:    text,
		IconUrl: iconUrl,
	}
	return e
}

func (e *Embed) SetImage(url string) *Embed {
	e.Image = EmbedImage{
		Url: url,
	}
	return e
}

func (e *Embed) SetThumbnail(url string) *Embed {
	e.Thumbnail = EmbedThumbnail{
		Url: url,
	}
	return e
}

func (e *Embed) SetVideo(url string) *Embed {
	e.Video = EmbedVideo{
		Url: url,
	}
	return e
}

func (e *Embed) SetProvider(name, url string) *Embed {
	e.Provider = EmbedProvider{
		Name: name,
		Url:  url,
	}
	return e
}

func (e *Embed) SetAuthor(name, url, iconUrl string) *Embed {
	e.Author = EmbedAuthor{
		Name:    name,
		Url:     url,
		IconUrl: iconUrl,
	}
	return e
}

func (e *Embed) AddField(title, content string, inline bool) *Embed {
	e.Fields = append(e.Fields, EmbedField{
		Name:   title,
		Value:  content,
		Inline: inline,
	})
	return e
}

func (e *Embed) AddBlankField(inline bool) *Embed {
	e.AddField("", "", inline)
	return e
}
