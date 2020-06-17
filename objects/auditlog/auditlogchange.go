package auditlog

// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key
type AuditLogChange struct {
	NewValue interface{} `json:"new_value"`
	OldValue interface{} `json:"old_value"`
	Key      ChangeKey   `json:"key"`
}
