package gateway

import (
	"context"
	"github.com/rxdn/gdl/gateway/payloads/events"
	"github.com/rxdn/gdl/objects/member"
	"github.com/sirupsen/logrus"
)

func RegisterCacheListeners(sm *ShardManager) {
	sm.RegisterListeners(
		readyListener,
		channelCreateListener,
		channelUpdateListener,
		channelDeleteListener,
		guildCreateListener,
		guildUpdateListener,
		guildDeleteListener,
		guildEmojisUpdateListeners,
		guildMemberAddListener,
		guildMemberRemoveListener,
		guildMemberUpdateListener,
		guildMembersChunkListener,
		guildRoleCreateListener,
		guildRoleUpdateListener,
		guildRoleDeleteListener,
		userUpdateListener,
		voiceStateUpdateListener,
	)
}

func readyListener(s *Shard, e *events.Ready) {
	logrus.Infof("shard %d: received ready", s.ShardId)

	s.sessionId = e.SessionId

	s.Cache.StoreSelf(context.Background(), e.User)
}

func channelCreateListener(s *Shard, e *events.ChannelCreate) {
	s.Cache.StoreChannel(context.Background(), e.Channel)
}

func channelUpdateListener(s *Shard, e *events.ChannelUpdate) {
	s.Cache.StoreChannel(context.Background(), e.Channel)
}

func channelDeleteListener(s *Shard, e *events.ChannelDelete) {
	s.Cache.DeleteChannel(context.Background(), e.Channel.Id)
}

func guildCreateListener(s *Shard, e *events.GuildCreate) {
	s.Cache.StoreGuild(context.Background(), e.Guild)
}

func guildUpdateListener(s *Shard, e *events.GuildUpdate) {
	s.Cache.StoreGuild(context.Background(), e.Guild)
}

func guildDeleteListener(s *Shard, e *events.GuildDelete) {
	s.Cache.DeleteGuild(context.Background(), e.Id)
}

func guildEmojisUpdateListeners(s *Shard, e *events.GuildEmojisUpdate) {
	for _, emoji := range e.Emojis {
		s.Cache.StoreEmoji(context.Background(), emoji, e.GuildId)
	}
}

func guildMemberAddListener(s *Shard, e *events.GuildMemberAdd) {
	s.Cache.StoreMember(context.Background(), e.Member, e.GuildId)
}

func guildMemberRemoveListener(s *Shard, e *events.GuildMemberRemove) {
	s.Cache.DeleteMember(context.Background(), e.User.Id, e.GuildId)
}

func guildMemberUpdateListener(s *Shard, e *events.GuildMemberUpdate) {
	s.Cache.StoreMember(context.Background(), member.Member{
		User:         e.User,
		Nick:         e.Nick,
		Roles:        e.Roles,
		PremiumSince: e.PremiumSince,
	}, e.GuildId)
}

func guildMembersChunkListener(s *Shard, e *events.GuildMembersChunk) {
	for _, member := range e.Members {
		s.Cache.StoreMember(context.Background(), member, e.GuildId)
	}
}

func guildRoleCreateListener(s *Shard, e *events.GuildRoleCreate) {
	s.Cache.StoreRole(context.Background(), e.Role, e.GuildId)
}

func guildRoleUpdateListener(s *Shard, e *events.GuildRoleUpdate) {
	s.Cache.StoreRole(context.Background(), e.Role, e.GuildId)
}

func guildRoleDeleteListener(s *Shard, e *events.GuildRoleDelete) {
	s.Cache.DeleteRole(context.Background(), e.RoleId)
}

func userUpdateListener(s *Shard, e *events.UserUpdate) {
	s.Cache.StoreUser(context.Background(), e.User)
}

func voiceStateUpdateListener(s *Shard, e *events.VoiceStateUpdate) {
	s.Cache.StoreVoiceState(context.Background(), e.VoiceState)
}
