package request

type RestError struct {
	ErrorCode int
	Message   string // json
}

func (r *RestError) IsClientError() bool {
	return r.ErrorCode >= 400 && r.ErrorCode < 500
}

func (r *RestError) IsServerError() bool {
	return r.ErrorCode >= 500
}

func (r RestError) Error() string {
	return r.Message
}
