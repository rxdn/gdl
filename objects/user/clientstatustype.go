package user

type ClientStatusType string

const (
	ClientStatusTypeOnline       ClientStatusType = "online"
	ClientStatusTypeIdle         ClientStatusType = "idle"
	ClientStatusTypeDoNotDisturb ClientStatusType = "dnd"
	ClientStatusTypeOffline      ClientStatusType = ""
)
