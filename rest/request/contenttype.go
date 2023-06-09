package request

type ContentType string

const (
	ApplicationJson           ContentType = "application/json"
	ApplicationFormUrlEncoded ContentType = "application/x-www-form-urlencoded"
	MultipartFormData         ContentType = "multipart/form-data"
	ImageJpeg                 ContentType = "image/jpeg"
	ImagePng                  ContentType = "image/png"
	ImageGif                  ContentType = "image/gif"
	Nil                       ContentType = ""
)
