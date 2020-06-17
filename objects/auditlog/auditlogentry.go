package auditlog

type AuditLogEntry struct {
	TargetId   uint64           `json:"target_id,string"`
	Changes    []AuditLogChange `json:"changes"`
	UserId     uint64           `json:"user_id,string"`
	Id         uint64           `json:"id,string"`
	ActionType AuditLogEvent    `json:"action_type"`
	Options    AuditEntryInfo   `json:"options"`
	Reason     *string          `json:"reason"`
}
