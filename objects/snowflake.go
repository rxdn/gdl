package objects

import (
	"github.com/rxdn/gdl/utils"
	"strconv"
)

type Snowflake uint64

func NewSnowflake(value uint64) Snowflake {
	return Snowflake(value)
}

func (s Snowflake) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strconv.FormatUint(uint64(s), 10) + "\""), nil
}

func (s *Snowflake) UnmarshalJSON(b []byte) error {
	parsed, err := utils.ReadStringUint64(b)
	if err != nil {
		return err
	}

	*s = Snowflake(parsed)

	return nil
}

func (s Snowflake) Value() uint64 {
	return uint64(s)
}
