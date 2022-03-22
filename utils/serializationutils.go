package utils

import (
	"fmt"
	"strconv"
)

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

func Contains[T comparable](slice []T, target T) bool {
	for _, el := range slice {
		if el == target {
			return true
		}
	}

	return false
}

func Initialise(m map[string]map[string]interface{}, key string) {
	if m[key] == nil {
		m[key] = make(map[string]interface{})
	}
}

func ReadStringUint16(s []byte) (uint16, error) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return 0, errMissingQuotation(s)
	}

	extracted := s[1 : len(s)-1]
	parsed, err := strconv.ParseUint(string(extracted), 10, 16)
	return uint16(parsed), err
}

func ReadStringUint64(s []byte) (uint64, error) {
	lower, upper := 0, len(s)

	if len(s) > 1 && s[lower] == '"' {
		lower++
	}

	if len(s) > 1 && s[upper-1] == '"' {
		upper--
	}

	extracted := s[lower:upper]
	parsed, err := strconv.ParseUint(string(extracted), 10, 64)
	return parsed, err
}

func errMissingQuotation(s []byte) error {
	return fmt.Errorf("string int is missing quotation marks: %s", string(s))
}
