package user

import (
	"fmt"
	"strconv"
)

type Avatar struct {
	Animated bool
	data []uint64
}

func (a *Avatar) String() string {
	if len(a.data) < 2 {
		return "" // blank avatar
	}

	var animatedPrefix string
	if a.Animated {
		animatedPrefix = "a_"
	}

	first := fmt.Sprintf("%016x", a.data[0])
	second := fmt.Sprintf("%016x", a.data[1])

	return fmt.Sprintf(`%s%s%s`, animatedPrefix, first, second)
}

func (a *Avatar) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, a.String())), nil
}

func (a *Avatar) UnmarshalJSON(data []byte) error {
	// 36 = animated
	// 34 = regular
	// != 36 & 34 = no avatar

	if len(data) == 36 {
		a.Animated = true
	} else if len(data) != 34 { // no avatar
		return nil
	}

	var first []byte
	var second []byte

	if a.Animated {
		first = data[3:19]
		second = data[19:35]
	} else {
		first = data[1:17]
		second = data[17:33]
	}

	a.data = make([]uint64, 2)

	var err error
	a.data[0], err = strconv.ParseUint(string(first), 16, 64)
	if err != nil {
		return err
	}

	a.data[1], err = strconv.ParseUint(string(second), 16, 64)
	if err != nil {
		return err
	}

	return nil
}


