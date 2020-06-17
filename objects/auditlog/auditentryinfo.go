package auditlog

// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditEntryInfo struct {
	DeleteMemberDays int        `json:"delete_member_days,string"` // MEMBER_PRUNE
	MembersRemoved   int        `json:"members_removed,string"`    // MEMBER_PRUNE
	ChannelId        uint64     `json:"channel_id,string"`         // MEMBER_MOVE & MESSAGE_PIN & MESSAGE_UNPIN & MESSAGE_DELETE
	MessageId        uint64     `json:"message_id,string"`         // MESSAGE_PIN & MESSAGE_UNPIN
	Id               uint64     `json:"id,string"`                 // MESSAGE_DELETE & MESSAGE_BULK_DELETE & MEMBER_DISCONNECT & MEMBER_MOVE
	Count            int        `json:"count,string"`              // CHANNEL_OVERWRITE_CREATE & CHANNEL_OVERWRITE_UPDATE & CHANNEL_OVERWRITE_DELETE
	Type             EntityType `json:"type"`                      // CHANNEL_OVERWRITE_CREATE & CHANNEL_OVERWRITE_UPDATE & CHANNEL_OVERWRITE_DELETE
	RoleName         string     `json:"role_name"`                 //CHANNEL_OVERWRITE_CREATE & CHANNEL_OVERWRITE_UPDATE & CHANNEL_OVERWRITE_DELETE
}
