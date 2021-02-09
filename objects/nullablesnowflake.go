package objects

import (
	"fmt"
	"github.com/rxdn/gdl/utils"
)

type NullableSnowflake uint64

func (i NullableSnowflake) MarshalJson() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", uint64(i))), nil
}

func (i *NullableSnowflake) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*i = 0
	} else {
		parsed, err := utils.ReadStringUint64(b)
		if err != nil {
			return err
		}

		*i = NullableSnowflake(parsed)
	}

	return nil
}