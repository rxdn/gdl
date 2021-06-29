package main

import (
	"encoding/json"
	"github.com/rxdn/gdl/rest/request"
	"testing"
)

func TestObjectError(t *testing.T) {
	var parsed request.ApiV8Error
	if err := json.Unmarshal(objectError, &parsed); err != nil {
		t.Error(err)
	}

	MustMatch(t, "code", parsed.Code, 50035)
	MustMatch(t, "message", parsed.Message, "Invalid Form Body")
	MustMatch(t, "errors: length", len(parsed.Errors), 1)


	err := parsed.Errors[0]
	MustMatch(t, "code", err.Code, "BASE_TYPE_REQUIRED")
	MustMatch(t, "message", err.Message, "This field is required")
	MustMatch(t, "fieldName", err.FieldName, "access_token")
}
