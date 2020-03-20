package utils

import "time"

func GetCurrentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
