package utils

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.Config{
	EscapeHTML:                    false,
	SortMapKeys:                   false,
	ObjectFieldMustBeSimpleString: false,
}.Froze()
