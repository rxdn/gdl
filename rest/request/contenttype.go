package request

type ContentType string

const (
	ApplicationJson           ContentType = "application/json"
	ApplicationFormUrlEncoded ContentType = "application/x-www-form-urlencoded"
	MultipartFormData         ContentType = "multipart/form-data"
	Nil                       ContentType = ""
)

type MultipartData interface {
	EncodeMultipartFormData() ([]byte, string, error)
}
