package request

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/constraints"
	"strconv"
	"strings"
)

type OAuthError struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e OAuthError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.ErrorDescription)
}

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
	Code    int
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
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := deNum(raw, "code", &e.Code); err != nil {
		return err
	}

	if err := de(raw, "message", &e.Message); err != nil {
		return err
	}

	var errors map[string]any
	if err := de(raw, "errors", &errors); err != nil {
		return err
	}

	for fieldName, value := range errors {
		var inner map[string]any
		if tmp, ok := value.(map[string]any); ok {
			inner = tmp
		} else {
			return fmt.Errorf("%s value is not an object", fieldName)
		}

		fieldErrors, err := deArray(inner, fieldName)
		if err != nil {
			return err
		}

		e.Errors = fieldErrors
	}

	return nil
}

func (e *ApiV8Error) FirstErrorCode() interface{} {
	if len(e.Errors) == 0 {
		return ""
	}

	return e.Errors[0].Code
}

func deArray(root map[string]any, fieldName string) ([]FieldError, error) {
	_, isObject := root["_errors"]
	if isObject {
		return deObj(root, fieldName)
	}

	var fieldErrors []FieldError
	for key, value := range root {
		var inner map[string]any
		if tmp, ok := value.(map[string]any); ok {
			inner = tmp
		} else {
			return nil, fmt.Errorf("%s value is not an object", key)
		}

		var innerFieldName string
		if _, err := strconv.Atoi(key); err == nil {
			innerFieldName = fmt.Sprintf("%s[%s]", fieldName, key)
		} else {
			innerFieldName = fmt.Sprintf("%s.%s", fieldName, key)
		}

		innerErrors, err := deArray(inner, innerFieldName)
		if err != nil {
			return nil, err
		}

		fieldErrors = append(fieldErrors, innerErrors...)
	}

	return fieldErrors, nil
}

func deObj(root map[string]any, fieldName string) ([]FieldError, error) {
	var errors []any
	if err := de(root, "_errors", &errors); err != nil {
		return nil, err
	}

	var fieldErrors []FieldError
	for _, err := range errors {
		var inner map[string]any
		if tmp, ok := err.(map[string]any); ok {
			inner = tmp
		} else {
			return nil, fmt.Errorf("%s _errors field value is not an object", fieldName)
		}

		var fieldError FieldError
		fieldError.FieldName = fieldName

		if err := de(inner, "code", &fieldError.Code); err != nil {
			return nil, err
		}

		if err := de(inner, "message", &fieldError.Message); err != nil {
			return nil, err
		}

		fieldErrors = append(fieldErrors, fieldError)
	}

	return fieldErrors, nil
}

func deNum[T constraints.Integer | constraints.Float](obj map[string]interface{}, key string, v *T) error {
	var f float64
	if err := de(obj, key, &f); err != nil {
		return err
	}

	*v = T(f)
	return nil
}

func de[T any](obj map[string]interface{}, key string, v *T) error {
	value, ok := obj[key]
	if !ok {
		return fmt.Errorf("%s was missing", key)
	}

	result, ok := value.(T)
	if !ok {
		return fmt.Errorf("%s was not a %T; got a %T", key, v, value)
	}

	*v = result
	return nil
}
