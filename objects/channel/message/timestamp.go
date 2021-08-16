package message

import (
	"fmt"
	"time"
)

type TimestampStyle string

const (
	TimestampStyleShortTime     TimestampStyle = "s"
	TimestampStyleLongTime      TimestampStyle = "T"
	TimestampStyleShortDate     TimestampStyle = "d"
	TimestampStyleLongDate      TimestampStyle = "D"
	TimestampStyleShortDateTime TimestampStyle = "f"
	TimestampStyleLongDateTime  TimestampStyle = "F"
	TimestampStyleRelativeTime  TimestampStyle = "R"
)

func BuildTimestamp(timestamp time.Time, style TimestampStyle) string {
	return BuildTimestampFromUnixSecs(timestamp.Unix(), style)
}

func BuildTimestampFromUnixSecs(timestamp int64, style TimestampStyle) string {
	return fmt.Sprintf("<t:%d:%s>", timestamp, style)
}