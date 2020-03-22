package request

type RequestType string

const (
	GET    RequestType = "GET"
	POST   RequestType = "POST"
	PATCH  RequestType = "PATCH"
	PUT    RequestType = "PUT"
	DELETE RequestType = "DELETE"
)
