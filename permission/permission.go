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
	UseApplicationCommands
	RequestToSpeak
	ManageEvents
	ManageThreads
	CreatePublicThreads
	CreatePrivateThreads
	UseExternalStickers
	SendMessagesInThreads
	UseEmbeddedActivities
	ModerateMembers
)

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

func (p Permission) String() string {
	switch p {
	case CreateInstantInvite:
		return "Create Instant Invite"
	case KickMembers:
		return "Kick Members"
	case BanMembers:
		return "Ban Members"
	case Administrator:
		return "Administrator"
	case ManageChannels:
		return "Manage Channels"
	case ManageGuild:
		return "Manage Server"
	case AddReactions:
		return "Add Reactions"
	case ViewAuditLog:
		return "View Audit Log"
	case PrioritySpeaker:
		return "Priority Speaker"
	case Stream:
		return "Stream"
	case ViewChannel:
		return "Read Messages"
	case SendMessages:
		return "Send Messages"
	case SendTTSMessages:
		return "Send TTS Messages"
	case ManageMessages:
		return "Manage Messages"
	case EmbedLinks:
		return "Embed Links"
	case AttachFiles:
		return "Attach Files"
	case ReadMessageHistory:
		return "Read Message History"
	case MentionEveryone:
		return "Mention Everyone"
	case UseExternalEmojis:
		return "Use External Emojis"
	case ViewGuildInsights:
		return "View Guild Insights"
	case Connect:
		return "Connect"
	case Speak:
		return "Speak"
	case MuteMembers:
		return "Mute Members"
	case DeafenMembers:
		return "Deafen Members"
	case MoveMembers:
		return "Move Members"
	case UseVAD:
		return "Use Voice Activity"
	case ChangeNickname:
		return "Change Nickname"
	case ManageNicknames:
		return "Manage Nicknames"
	case ManageRoles:
		return "Manage Roles"
	case ManageWebhooks:
		return "Manage Webhooks"
	case ManageEmojis:
		return "Manage Emojis"
	case UseApplicationCommands:
		return "Use Application Commands"
	case RequestToSpeak:
		return "Request to Speak"
	case ManageEvents:
		return "Manage Events"
	case ManageThreads:
		return "Manage Threads"
	case CreatePublicThreads:
		return "Create Public Threads"
	case CreatePrivateThreads:
		return "Create Private Threads"
	case UseExternalStickers:
		return "Use External Stickers"
	case SendMessagesInThreads:
		return "Send Messages in Threads"
	case UseEmbeddedActivities:
		return "Use Embedded Activities"
	case ModerateMembers:
		return "Moderate Members"
	default:
		return "Unknown Permission"
	}
}

var AllPermissions = []Permission{
	CreateInstantInvite,
	KickMembers,
	BanMembers,
	Administrator,
	ManageChannels,
	ManageGuild,
	AddReactions,
	ViewAuditLog,
	PrioritySpeaker,
	Stream,
	ViewChannel,
	SendMessages,
	SendTTSMessages,
	ManageMessages,
	EmbedLinks,
	AttachFiles,
	ReadMessageHistory,
	MentionEveryone,
	UseExternalEmojis,
	ViewGuildInsights,
	Connect,
	Speak,
	MuteMembers,
	DeafenMembers,
	MoveMembers,
	UseVAD,
	ChangeNickname,
	ManageNicknames,
	ManageRoles,
	ManageWebhooks,
	ManageEmojis,
	UseApplicationCommands,
	RequestToSpeak,
	ManageEvents,
	ManageThreads,
	CreatePublicThreads,
	CreatePrivateThreads,
	UseExternalStickers,
	SendMessagesInThreads,
	UseEmbeddedActivities,
	ModerateMembers,
}
