package utils

import (
	"fmt"
	"github.com/fatih/structs"
	"reflect"
	"strconv"
)

func Unmarshal(raw []byte, v interface{}) {
	_ = json.Unmarshal(raw, &v)
}

func AppendElem(m map[string]map[string]interface{}, key string, elem map[string]interface{}) map[string]map[string]interface{} {
	for k, v := range elem {
		m[key][k] = v
	}
	return m
}

func Append(m, elem map[string]map[string]interface{}) map[string]map[string]interface{} {
	for k := range elem {
		Initialise(m, k)
		for k1, v1 := range elem[k] {
			m[k][k1] = v1
		}
	}
	return m
}

func Contains(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {

			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}

	return false
}

func Initialise(m map[string]map[string]interface{}, key string) {
	if m[key] == nil {
		m[key] = make(map[string]interface{})
	}
}

func CopyNonNil(m map[string]map[string]interface{}, keyName string, obj interface{}) {
	Initialise(m, keyName)

	for k, v := range structs.Map(obj) {
		if !IsZero(reflect.ValueOf(v)) {
			m[keyName][k] = v
		}
	}
}

func ReadStringUint16(s []byte) (uint16, error) {
	if len(s) < 2 || s[0] != '"' || s[len(s) - 1] != '"' {
		return 0, errMissingQuotation(s)
	}

	extracted := s[1:len(s)-1]
	parsed, err := strconv.ParseUint(string(extracted), 10, 16)
	return uint16(parsed), err
}

func ReadStringUint64(s []byte) (uint64, error) {
	if len(s) < 2 || s[0] != '"' || s[len(s) - 1] != '"' {
		return 0, errMissingQuotation(s)
	}

	extracted := s[1:len(s)-1]
	parsed, err := strconv.ParseUint(string(extracted), 10, 64)
	return parsed, err
}

func errMissingQuotation(s []byte) error {
	return fmt.Errorf("string int is missing quotation marks: %s", string(s))
}