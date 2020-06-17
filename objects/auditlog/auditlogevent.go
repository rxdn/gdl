package auditlog

type AuditLogEvent int

const (
	EventGuildUpdate              AuditLogEvent = 1
	EventChannelCreate            AuditLogEvent = 10
	EventChannelUpdate            AuditLogEvent = 11
	EventChannelDelete            AuditLogEvent = 12
	EventChannelOverwriteCreate   AuditLogEvent = 13
	EventChannelOverwriteUpdate   AuditLogEvent = 14
	EventChannelOverwriteDelete   AuditLogEvent = 15
	EventMemberKick               AuditLogEvent = 20
	EventMemberPrune              AuditLogEvent = 21
	EventMemberBanAdd             AuditLogEvent = 22
	EventMemberBanRemove          AuditLogEvent = 23
	EventMemberUpdate             AuditLogEvent = 24
	EventMemberRoleUpdate         AuditLogEvent = 25
	EventMemberMove               AuditLogEvent = 26
	EventMemberDisconnect         AuditLogEvent = 27
	EventBotAdd                   AuditLogEvent = 28
	EventRoleCreate               AuditLogEvent = 30
	EventRoleUpdate               AuditLogEvent = 31
	EventRoleDelete               AuditLogEvent = 32
	EventInviteCreate             AuditLogEvent = 40
	EventInviteUpdate             AuditLogEvent = 41
	EventInviteDelete             AuditLogEvent = 42
	EventWebhookCreate            AuditLogEvent = 50
	EventWebhookUpdate            AuditLogEvent = 51
	EventWebhookDelete            AuditLogEvent = 52
	EventEmojiCreate              AuditLogEvent = 60
	EventEmojiUpdate              AuditLogEvent = 61
	EventEmojiDelete              AuditLogEvent = 62
	EventMessageDelete            AuditLogEvent = 72
	EventMessageBulkDelete        AuditLogEvent = 73
	EventMessagePin               AuditLogEvent = 74
	EventMessageUnpin             AuditLogEvent = 75
	EventMessageIntegrationCreate AuditLogEvent = 80
	EventMessageIntegrationUpdate AuditLogEvent = 81
	EventMessageIntegrationDelete AuditLogEvent = 82
)
