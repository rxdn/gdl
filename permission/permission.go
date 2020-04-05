package permission

type Permission int

const (
	CreateInstantInvite Permission = 0x00000001
	KickMembers         Permission = 0x00000002
	BanMembers          Permission = 0x00000004
	Administrator       Permission = 0x00000008
	ManageChannels      Permission = 0x00000010
	ManageGuild         Permission = 0x00000020
	AddReactions        Permission = 0x00000040
	ViewAuditLog        Permission = 0x00000080
	ViewChannel         Permission = 0x00000400 // Read messages
	SendMessages        Permission = 0x00000800
	SendTTSMessages     Permission = 0x00001000
	ManageMessages      Permission = 0x00002000
	EmbedLinks          Permission = 0x00004000
	AttachFiles         Permission = 0x00008000
	ReadMessageHistory  Permission = 0x00010000
	MentionEveryone     Permission = 0x00020000
	Connect             Permission = 0x00100000
	Speak               Permission = 0x00200000
	MuteMembers         Permission = 0x00400000
	DeafenMembers       Permission = 0x00800000
	MoveMembers         Permission = 0x01000000
	UseVAD              Permission = 0x02000000 // Use voice activity
	PrioritySpeaker     Permission = 0x00000100
	ChangeNickname      Permission = 0x04000000
	ManageNicknames     Permission = 0x08000000
	ManageRoles         Permission = 0x10000000 // Manage permissions
	ManageWebhooks      Permission = 0x20000000
	ManageEmojis        Permission = 0x40000000
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

func HasPermissionRaw(permissions int, permission Permission) bool {
	return permissions&int(permission) == int(permission)
}

func BuildPermissions(permissions ...Permission) int {
	i := 0

	for _, permission := range permissions {
		i |= int(permission)
	}

	return i
}

func twosComplement(i int) int {
	return (-i) - 1
}
