package objects

import (
	"fmt"
	"github.com/rxdn/gdl/utils"
	"strconv"
)

type NullableSnowflake struct {
	IsNull bool
	Value  uint64
}

func NewNullableSnowflake(value uint64) NullableSnowflake {
	return NullableSnowflake{
		IsNull: false,
		Value:  value,
	}
}

func NewNullSnowflake() NullableSnowflake {
	return NullableSnowflake{
		IsNull: true,
		Value:  0,
	}
}

func (i NullableSnowflake) MarshalJson() ([]byte, error) {
	if i.IsNull {
		return []byte("null"), nil
	} else {
		return []byte("\"" + strconv.FormatUint(i.Value, 10) + "\""), nil
	}
}

func (i *NullableSnowflake) UnmarshalJSON(b []byte) error {
	*i = NewNullSnowflake()

	if string(b) == "null" {
		i.IsNull = true
	} else {
		parsed, err := utils.ReadStringUint64(b)
		if err != nil {
			return err
		}

		i.IsNull = false
		i.Value = parsed
	}

	return nil
}
