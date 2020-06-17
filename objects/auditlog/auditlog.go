package auditlog

import (
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/integration"
	"github.com/rxdn/gdl/objects/user"
)

type AuditLog struct {
	Webhooks     []guild.Webhook           `json:"webhooks"`
	Users        []user.User               `json:"users"`
	Entries      []AuditLogEntry           `json:"audit_log_entries"`
	Integrations []integration.Integration `json:"integrations"`
}
