package request

import (
	"fmt"
	"strings"
)

type RestError struct {
	StatusCode int
	ApiError   ApiV8Error
	Url        string
	Raw        []byte
}

type ApiV8Error struct {
	Code   interface{} `json:"code"` // Can be int or string
	Errors struct {
		Errors []ApiV8Error `json:"_errors"`
	} `json:"errors"`
	Message string `json:"message"`
}

func (r *RestError) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

func (r *RestError) IsServerError() bool {
	return r.StatusCode >= 500
}

func (r RestError) Error() string {
	if len(r.ApiError.Errors.Errors) > 0 {
		messages := make([]string, len(r.ApiError.Errors.Errors))
		for i, err := range r.ApiError.Errors.Errors {
			messages[i] = err.Message
		}

		return fmt.Sprintf("%s: %s", r.ApiError.Message, strings.Join(messages, ", "))
	} else {
		return r.ApiError.Message
	}
}
