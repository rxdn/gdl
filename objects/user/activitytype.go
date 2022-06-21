package user

type ActivityType int

const (
	ActivityTypePlaying ActivityType = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeWatching
	ActivityTypeCustom
	ActivityTypeCompeting
)
