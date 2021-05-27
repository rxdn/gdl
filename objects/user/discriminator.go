package user

import (
	"fmt"
	"github.com/rxdn/gdl/utils"
)

// Jsoniter does not like leading zeroes on ints
type Discriminator uint16

func (d Discriminator) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", uint16(d))), nil
}

func (d *Discriminator) UnmarshalJSON(b []byte) error {
	parsed, err := utils.ReadStringUint16(b)
	if err != nil {
		return err
	}

	*d = Discriminator(parsed)
	return nil
}
