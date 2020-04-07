package guild

import (
	"github.com/rxdn/gdl/objects/member"
)

type CachedVoiceState struct {
	ChannelId uint64 `db:"channel_id"`
	SessionId string `db:"session_id"`
	Deaf      bool   `db:"deaf"`
	Mute      bool   `db:"mute"`
	SelfDeaf  bool   `db:"self_deaf"`
	SelfMute  bool   `db:"self_mute"`
	Suppress  bool   `db:"suppress"`
}

func (s *CachedVoiceState) ToVoiceState(guildId uint64, m member.Member) VoiceState {
	return VoiceState{
		GuildId:   guildId,
		ChannelId: s.ChannelId,
		UserId:    m.User.Id,
		Member:    m,
		SessionId: s.SessionId,
		Deaf:      s.Deaf,
		Mute:      s.Mute,
		SelfDeaf:  s.SelfDeaf,
		SelfMute:  s.SelfMute,
		Suppress:  s.Suppress,
	}
}
