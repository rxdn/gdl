package guild

import "github.com/rxdn/gdl/objects/member"

type VoiceState struct {
	GuildId   uint64        `json:"guild_id,string"`
	ChannelId uint64        `json:"channel_id,string"`
	UserId    uint64        `json:"user_id,string"`
	Member    member.Member `json:"member"`
	SessionId string        `json:"session_id"`
	Deaf      bool          `json:"deaf"`
	Mute      bool          `json:"mute"`
	SelfDeaf  bool          `json:"self_deaf"`
	SelfMute  bool          `json:"self_mute"`
	Suppress  bool          `json:"suppress"`
}

func (s *VoiceState) ToCachedVoiceState() CachedVoiceState {
	return CachedVoiceState{
		ChannelId: s.ChannelId,
		SessionId: s.SessionId,
		Deaf:      s.Deaf,
		Mute:      s.Mute,
		SelfDeaf:  s.SelfDeaf,
		SelfMute:  s.SelfMute,
		Suppress:  s.Suppress,
	}
}
