package interaction

type ApplicationCommand struct {
	Id                uint64                     `json:"id,string,omitempty"`
	Type              ApplicationCommandType     `json:"application_command_type"`
	ApplicationId     uint64                     `json:"application_id,string,omitempty"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Options           []ApplicationCommandOption `json:"options"`
	DefaultPermission bool                       `json:"default_permission,omitempty"`
}

type ApplicationCommandType uint8

const (
	ApplicationCommandTypeChatInput ApplicationCommandType = iota + 1
	ApplicationCommandTypeUser
	ApplicationCommandTypeMessage
)
