package request

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RestError struct {
	StatusCode int
	ApiError   ApiV8Error
	Url        string
	Raw        []byte
}

func (r *RestError) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

func (r *RestError) IsServerError() bool {
	return r.StatusCode >= 500
}

func (r RestError) Error() string {
	if len(r.ApiError.Errors) > 0 {
		messages := make([]string, len(r.ApiError.Errors))
		for i, err := range r.ApiError.Errors {
			messages[i] = err.Error()
		}

		return fmt.Sprintf("%s: %s", r.ApiError.Message, strings.Join(messages, ", "))
	} else {
		return r.ApiError.Message
	}
}

type ApiV8Error struct {
	Code    interface{} // Can be int or string
	Errors  []FieldError
	Message string
}

type FieldError struct {
	FieldName string      `json:"-"`
	Code      interface{} `json:"code"` // Can be int or string
	Message   string      `json:"message"`
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %s (%v)", e.FieldName, e.Message, e.Code)
}

func (e *ApiV8Error) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.Code = raw["code"]
	if f, ok := e.Code.(float64); ok {
		e.Code = int(f)
	}

	var ok bool
	e.Message, ok = raw["message"].(string)
	if !ok {
		return fmt.Errorf("message was not a string")
	}

	errors, ok := raw["errors"].(map[string]interface{})
	if ok { // Only 400s have this field
		for fieldName, value := range errors {
			value, ok := value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%s value is not an object", fieldName)
			}

			errors, ok := value["_errors"].([]interface{})
			if ok { // Is object error
				for _, err := range errors {
					err, ok := err.(map[string]interface{})
					if !ok {
						return fmt.Errorf("%s _errors field value is not an object", fieldName)
					}

					e.Errors = append(e.Errors, deFieldError(fieldName, err))
				}
			} else { // Is array error
				for i, entry := range value {
					value, ok := entry.(map[string]interface{})
					if !ok {
						return fmt.Errorf("%s value is not an object", fieldName)
					}

					for fieldName, value := range value {
						value, ok := value.(map[string]interface{})
						if !ok {
							return fmt.Errorf("%s value is not an object", fieldName)
						}

						errors, ok := value["_errors"].([]interface{})
						if !ok {
							return fmt.Errorf("%s array entry %s: _errors field is not an array", fieldName, i)
						}

						for _, err := range errors {
							err, ok := err.(map[string]interface{})
							if !ok {
								return fmt.Errorf("%s _errors field value is not an object", fieldName)
							}

							e.Errors = append(e.Errors, deFieldError(fieldName, err))
						}
					}
				}
			}
		}
	}

	return nil
}

func (e *ApiV8Error) FirstErrorCode() interface{} {
	if len(e.Errors) == 0 {
		return ""
	}

	return e.Errors[0].Code
}

func deFieldError(fieldName string, data map[string]interface{}) (err FieldError) {
	err.FieldName = fieldName
	err.Code = data["code"]
	err.Message, _ = data["message"].(string)

	if f, ok := err.Code.(float64); ok {
		err.Code = int(f)
	}

	return
}
