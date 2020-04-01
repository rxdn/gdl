package channel

type ChannelType int

const (
	ChannelTypeGuildText     ChannelType = 0
	ChannelTypeDM            ChannelType = 1
	ChannelTypeGuildVoice    ChannelType = 2
	ChannelTypeGroupDM       ChannelType = 3
	ChannelTypeGuildCategory ChannelType = 4
	ChannelTypeGuildNews     ChannelType = 5
	ChannelTypeGuildStore    ChannelType = 6
)
