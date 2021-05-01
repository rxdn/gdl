package channel

import "time"

type ThreadMetadata struct {
	Archived            bool      `json:"archived"`
	ArchiverId          uint64    `json:"archiver_id,string"`
	AutoArchiveDuration uint16    `json:"auto_archive_duration"`
	ArchiveTimestamp    time.Time `json:"archive_timestamp"`
	Locked              bool      `json:"locked"`
}

type ThreadMember struct {
	ThreadId      uint64    `json:"id,string"`
	UserId        uint64    `json:"user_id,string"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	Flags         uint      `json:"flags"`
}
