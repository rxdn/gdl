package user

type ClientStatusType string

const (
	ClientStatusTypeOnline       ClientStatusType = "online"
	ClientStatusTypeIdle         ClientStatusType = "idle"
	ClientStatusTypeDoNotDisturb ClientStatusType = "dnd"
	ClientStatusTypeInvisible    ClientStatusType = "invisible"
	ClientStatusTypeOffline      ClientStatusType = "offline"
)
