package permission

type Permission uint64

const (
	CreateInstantInvite Permission = 1 << iota
	KickMembers
	BanMembers
	Administrator
	ManageChannels
	ManageGuild
	AddReactions
	ViewAuditLog
	PrioritySpeaker
	Stream
	ViewChannel // Read messages
	SendMessages
	SendTTSMessages
	ManageMessages
	EmbedLinks
	AttachFiles
	ReadMessageHistory
	MentionEveryone
	UseExternalEmojis
	ViewGuildInsights
	Connect
	Speak
	MuteMembers
	DeafenMembers
	MoveMembers
	UseVAD // Use voice activity
	ChangeNickname
	ManageNicknames
	ManageRoles // Manage permissions
	ManageWebhooks
	ManageEmojis
	UseSlashCommands
	RequestToSpeak
	ManageThreads
	UsePublicThreads
	UsePrivateThreads
)

var AllPermissions = []Permission{
	CreateInstantInvite,
	KickMembers,
	BanMembers,
	Administrator,
	ManageChannels,
	ManageGuild,
	AddReactions,
	ViewAuditLog,
	ViewAuditLog,
	ViewChannel,
	SendMessages,
	SendTTSMessages,
	ManageMessages,
	EmbedLinks,
	AttachFiles,
	ReadMessageHistory,
	MentionEveryone,
	Connect,
	Speak,
	MuteMembers,
	DeafenMembers,
	MoveMembers,
	UseVAD,
	PrioritySpeaker,
	ChangeNickname,
	ManageNicknames,
	ManageRoles,
	ManageWebhooks,
	ManageEmojis,
}

func HasPermissionRaw(permissions uint64, permission Permission) bool {
	return permissions&uint64(permission) == uint64(permission)
}

func BuildPermissions(permissions ...Permission) uint64 {
	var i uint64

	for _, permission := range permissions {
		i |= uint64(permission)
	}

	return i
}
