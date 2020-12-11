package interaction

type ApplicationCommandOptionType uint8

const (
	OptionTypeSubCommand ApplicationCommandOptionType = iota + 1
	OptionTypeSubCommandGroup
	OptionTypeString
	OptionTypeInteger
	OptionTypeBoolean
	OptionTypeUser
	OptionTypeChannel
	OptionTypeRole
)
