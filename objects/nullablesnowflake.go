package objects

import (
	"fmt"
	"github.com/rxdn/gdl/utils"
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
	return []byte(fmt.Sprintf("\"%d\"", i.Value)), nil
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
