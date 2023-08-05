package gateway

import (
	"context"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/auditlog"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/integration"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/objects/invite"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest"
)

func (s *Shard) GetChannel(ctx context.Context, channelId uint64) (channel.Channel, error) {
	if s.Cache.Options().Channels {
		if cached, err := s.Cache.GetChannel(ctx, channelId); err == nil {
			return cached, nil
		} else if err != cache.ErrNotFound {
			return channel.Channel{}, err
		}
	}

	ch, err := rest.GetChannel(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
	if err != nil {
		return channel.Channel{}, err
	}

	if s.Cache.Options().Channels {
		if err := s.Cache.StoreChannel(ctx, ch); err != nil {
			return channel.Channel{}, err
		}
	}

	return ch, err
}

func (s *Shard) ModifyChannel(ctx context.Context, channelId uint64, data rest.ModifyChannelData) (channel.Channel, error) {
	ch, err := rest.ModifyChannel(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
	if err != nil {
		return channel.Channel{}, err
	}

	if s.Cache.Options().Channels {
		if err := s.Cache.StoreChannel(ctx, ch); err != nil {
			return channel.Channel{}, err
		}
	}

	return ch, err
}

func (s *Shard) DeleteChannel(ctx context.Context, channelId uint64) (channel.Channel, error) {
	return rest.DeleteChannel(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) GetChannelMessages(ctx context.Context, channelId uint64, options rest.GetChannelMessagesData) ([]message.Message, error) {
	return rest.GetChannelMessages(ctx, s.Token, s.ShardManager.RateLimiter, channelId, options)
}

func (s *Shard) GetChannelMessage(ctx context.Context, channelId, messageId uint64) (message.Message, error) {
	return rest.GetChannelMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) CreateMessage(ctx context.Context, channelId uint64, content string) (message.Message, error) {
	return s.CreateMessageComplex(ctx, channelId, rest.CreateMessageData{
		Content: content,
	})
}

func (s *Shard) CreateMessageReply(ctx context.Context, channelId uint64, content string, reference *message.MessageReference) (message.Message, error) {
	return s.CreateMessageComplex(ctx, channelId, rest.CreateMessageData{
		Content:          content,
		MessageReference: reference,
	})
}

func (s *Shard) CreateMessageEmbed(ctx context.Context, channelId uint64, embed ...*embed.Embed) (message.Message, error) {
	return s.CreateMessageComplex(ctx, channelId, rest.CreateMessageData{
		Embeds: embed,
	})
}

func (s *Shard) CreateMessageEmbedReply(ctx context.Context, channelId uint64, e *embed.Embed, reference *message.MessageReference) (message.Message, error) {
	return s.CreateMessageComplex(ctx, channelId, rest.CreateMessageData{
		Embeds:           []*embed.Embed{e},
		MessageReference: reference,
	})
}

func (s *Shard) CreateMessageComplex(ctx context.Context, channelId uint64, data rest.CreateMessageData) (message.Message, error) {
	return rest.CreateMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) CreateReaction(ctx context.Context, channelId, messageId uint64, emoji string) error {
	return rest.CreateReaction(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId, emoji)
}

func (s *Shard) DeleteOwnReaction(ctx context.Context, channelId, messageId uint64, emoji string) error {
	return rest.DeleteOwnReaction(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId, emoji)
}

func (s *Shard) DeleteUserReaction(ctx context.Context, channelId, messageId, userId uint64, emoji string) error {
	return rest.DeleteUserReaction(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId, userId, emoji)
}

func (s *Shard) GetReactions(ctx context.Context, channelId, messageId uint64, emoji string, options rest.GetReactionsData) ([]user.User, error) {
	return rest.GetReactions(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId, emoji, options)
}

func (s *Shard) DeleteAllReactions(ctx context.Context, channelId, messageId uint64) error {
	return rest.DeleteAllReactions(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) DeleteAllReactionsEmoji(ctx context.Context, channelId, messageId uint64, emoji string) error {
	return rest.DeleteAllReactionsEmoji(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId, emoji)
}

func (s *Shard) EditMessage(ctx context.Context, channelId, messageId uint64, data rest.EditMessageData) (message.Message, error) {
	return rest.EditMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId, data)
}

func (s *Shard) DeleteMessage(ctx context.Context, channelId, messageId uint64) error {
	return rest.DeleteMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) BulkDeleteMessages(ctx context.Context, channelId uint64, messages []uint64) error {
	return rest.BulkDeleteMessages(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messages)
}

func (s *Shard) EditChannelPermissions(ctx context.Context, channelId uint64, updated channel.PermissionOverwrite) error {
	return rest.EditChannelPermissions(ctx, s.Token, s.ShardManager.RateLimiter, channelId, updated)
}

func (s *Shard) GetChannelInvites(ctx context.Context, channelId uint64) ([]invite.InviteMetadata, error) {
	return rest.GetChannelInvites(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) CreateChannelInvite(ctx context.Context, channelId uint64, data rest.CreateInviteData) (invite.Invite, error) {
	return rest.CreateChannelInvite(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) DeleteChannelPermissions(ctx context.Context, channelId, overwriteId uint64) error {
	return rest.DeleteChannelPermissions(ctx, s.Token, s.ShardManager.RateLimiter, channelId, overwriteId)
}

func (s *Shard) TriggerTypingIndicator(ctx context.Context, channelId uint64) error {
	return rest.TriggerTypingIndicator(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) GetPinnedMessages(ctx context.Context, channelId uint64) ([]message.Message, error) {
	return rest.GetPinnedMessages(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) AddPinnedChannelMessage(ctx context.Context, channelId, messageId uint64) error {
	return rest.AddPinnedChannelMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) DeletePinnedChannelMessage(ctx context.Context, channelId, messageId uint64) error {
	return rest.DeletePinnedChannelMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) JoinThread(ctx context.Context, channelId uint64) error {
	return rest.JoinThread(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) AddThreadMember(ctx context.Context, channelId, userId uint64) error {
	return rest.AddThreadMember(ctx, s.Token, s.ShardManager.RateLimiter, channelId, userId)
}

func (s *Shard) LeaveThread(ctx context.Context, channelId uint64) error {
	return rest.LeaveThread(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) RemoveThreadMember(ctx context.Context, channelId, userId uint64) error {
	return rest.RemoveThreadMember(ctx, s.Token, s.ShardManager.RateLimiter, channelId, userId)
}

func (s *Shard) GetThreadMember(ctx context.Context, channelId, userId uint64) (channel.ThreadMember, error) {
	return rest.GetThreadMember(ctx, s.Token, s.ShardManager.RateLimiter, channelId, userId)
}

func (s *Shard) ListThreadMembers(ctx context.Context, channelId uint64) ([]channel.ThreadMember, error) {
	return rest.ListThreadMembers(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) ListActiveThreads(ctx context.Context, guildId uint64) (rest.ThreadsResponse, error) {
	return rest.ListActiveThreads(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) ListPublicArchivedThreads(ctx context.Context, channelId uint64, data rest.ListThreadsData) (rest.ThreadsResponse, error) {
	return rest.ListPublicArchivedThreads(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) ListPrivateArchivedThreads(ctx context.Context, channelId uint64, data rest.ListThreadsData) (rest.ThreadsResponse, error) {
	return rest.ListPrivateArchivedThreads(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) ListJoinedPrivateArchivedThreads(ctx context.Context, channelId uint64, data rest.ListThreadsData) (rest.ThreadsResponse, error) {
	return rest.ListPrivateArchivedThreads(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) StartThreadWithMessage(ctx context.Context, channelId, messageId uint64, data rest.StartThreadWithMessageData) (channel.Channel, error) {
	return rest.StartThreadWithMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, messageId, data)
}

func (s *Shard) StartThreadWithoutMessage(ctx context.Context, channelId, messageId uint64, data rest.StartThreadWithoutMessageData) (channel.Channel, error) {
	return rest.StartThreadWithoutMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) CreatePublicThread(ctx context.Context, channelId uint64, name string, autoArchiveDuration uint16) (channel.Channel, error) {
	data := rest.StartThreadWithoutMessageData{
		Name:                name,
		AutoArchiveDuration: autoArchiveDuration,
		Type:                channel.ChannelTypeGuildPublicThread,
	}

	return rest.StartThreadWithoutMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) CreatePrivateThread(ctx context.Context, channelId uint64, name string, autoArchiveDuration uint16, invitable bool) (channel.Channel, error) {
	data := rest.StartThreadWithoutMessageData{
		Name:                name,
		AutoArchiveDuration: autoArchiveDuration,
		Type:                channel.ChannelTypeGuildPrivateThread,
		Invitable:           invitable,
	}

	return rest.StartThreadWithoutMessage(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) ListGuildEmojis(ctx context.Context, guildId uint64) ([]emoji.Emoji, error) {
	if s.Cache.Options().Emojis && s.Cache.Options().Guilds {
		if emojis, err := s.Cache.GetGuildEmojis(ctx, guildId); err == nil {
			return emojis, nil
		} else if err != cache.ErrNotFound {
			return nil, err
		}
	}

	emojis, err := rest.ListGuildEmojis(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
	if err != nil {
		return nil, err
	}

	if s.Cache.Options().Emojis {
		if err := s.Cache.StoreEmojis(ctx, emojis, guildId); err != nil {
			return nil, err
		}
	}

	return emojis, err
}

func (s *Shard) GetGuildEmoji(ctx context.Context, guildId uint64, emojiId uint64) (emoji.Emoji, error) {
	if s.Cache.Options().Emojis {
		if e, err := s.Cache.GetEmoji(ctx, emojiId); err == nil {
			return e, nil
		} else if err != cache.ErrNotFound {
			return emoji.Emoji{}, err
		}
	}

	e, err := rest.GetGuildEmoji(ctx, s.Token, s.ShardManager.RateLimiter, guildId, emojiId)
	if err != nil {
		return emoji.Emoji{}, err
	}

	if s.Cache.Options().Emojis {
		if err := s.Cache.StoreEmoji(ctx, e, guildId); err != nil {
			return emoji.Emoji{}, err
		}
	}

	return e, err
}

func (s *Shard) CreateGuildEmoji(ctx context.Context, guildId uint64, data rest.CreateEmojiData) (emoji.Emoji, error) {
	return rest.CreateGuildEmoji(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) ModifyGuildEmoji(ctx context.Context, guildId, emojiId uint64, data rest.CreateEmojiData) (emoji.Emoji, error) {
	return rest.ModifyGuildEmoji(ctx, s.Token, s.ShardManager.RateLimiter, guildId, emojiId, data)
}

func (s *Shard) CreateGuild(ctx context.Context, data rest.CreateGuildData) (guild.Guild, error) {
	return rest.CreateGuild(ctx, s.Token, data)
}

func (s *Shard) GetGuild(ctx context.Context, guildId uint64) (guild.Guild, error) {
	if s.Cache.Options().Guilds {
		if cached, err := s.Cache.GetGuild(ctx, guildId); err == nil {
			return cached, nil
		} else if err != cache.ErrNotFound {
			return guild.Guild{}, err
		}
	}

	g, err := rest.GetGuild(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
	if err != nil {
		return guild.Guild{}, err
	}

	if s.Cache.Options().Guilds {
		if err := s.Cache.StoreGuild(ctx, g); err != nil {
			return guild.Guild{}, err
		}
	}

	return g, err
}

func (s *Shard) GetGuildPreview(ctx context.Context, guildId uint64) (guild.GuildPreview, error) {
	return rest.GetGuildPreview(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) ModifyGuild(ctx context.Context, guildId uint64, data rest.ModifyGuildData) (guild.Guild, error) {
	return rest.ModifyGuild(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) DeleteGuild(ctx context.Context, guildId uint64) error {
	return rest.DeleteGuild(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetGuildChannels(ctx context.Context, guildId uint64) ([]channel.Channel, error) {
	if s.Cache.Options().Channels {
		if cached, err := s.Cache.GetGuildChannels(ctx, guildId); err == nil {
			return cached, nil
		} else if err != cache.ErrNotFound {
			return nil, err
		}
	}

	channels, err := rest.GetGuildChannels(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
	if err != nil {
		return nil, err
	}

	if s.Cache.Options().Channels {
		if err := s.Cache.StoreChannels(ctx, channels); err != nil {
			return nil, err
		}
	}

	return channels, err
}

func (s *Shard) CreateGuildChannel(ctx context.Context, guildId uint64, data rest.CreateChannelData) (channel.Channel, error) {
	return rest.CreateGuildChannel(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) ModifyGuildChannelPositions(ctx context.Context, guildId uint64, positions []rest.Position) error {
	return rest.ModifyGuildChannelPositions(ctx, s.Token, s.ShardManager.RateLimiter, guildId, positions)
}

func (s *Shard) GetGuildMember(ctx context.Context, guildId, userId uint64) (member.Member, error) {
	if s.Cache.Options().Members {
		if cached, err := s.Cache.GetMember(ctx, guildId, userId); err == nil {
			return cached, nil
		} else if err != cache.ErrNotFound {
			return member.Member{}, err
		}
	}

	m, err := rest.GetGuildMember(ctx, s.Token, s.ShardManager.RateLimiter, guildId, userId)
	if err != nil {
		return member.Member{}, err
	}

	if s.Cache.Options().Members {
		if err := s.Cache.StoreMember(ctx, m, guildId); err != nil {
			return member.Member{}, err
		}
	}

	return m, err
}

func (s *Shard) SearchGuildMembers(ctx context.Context, guildId uint64, data rest.SearchGuildMembersData) ([]member.Member, error) {
	members, err := rest.SearchGuildMembers(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
	if err != nil {
		return nil, err
	}

	if s.Cache.Options().Members {
		if err := s.Cache.StoreMembers(ctx, members, guildId); err != nil {
			return nil, err
		}

		var users []user.User
		for _, m := range members {
			users = append(users, m.User)
		}

		if err := s.Cache.StoreUsers(ctx, users); err != nil {
			return nil, err
		}
	}

	return members, err
}

func (s *Shard) ListGuildMembers(ctx context.Context, guildId uint64, data rest.ListGuildMembersData) ([]member.Member, error) {
	members, err := rest.ListGuildMembers(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
	if err != nil {
		return nil, err
	}

	if s.Cache.Options().Members {
		if err := s.Cache.StoreMembers(ctx, members, guildId); err != nil {
			return nil, err
		}

		var users []user.User
		for _, m := range members {
			users = append(users, m.User)
		}

		if err := s.Cache.StoreUsers(ctx, users); err != nil {
			return nil, err
		}
	}

	return members, err
}

func (s *Shard) ModifyGuildMember(ctx context.Context, guildId, userId uint64, data rest.ModifyGuildMemberData) error {
	return rest.ModifyGuildMember(ctx, s.Token, s.ShardManager.RateLimiter, guildId, userId, data)
}

func (s *Shard) ModifyCurrentUserNick(ctx context.Context, guildId uint64, nick string) error {
	return rest.ModifyCurrentUserNick(ctx, s.Token, s.ShardManager.RateLimiter, guildId, nick)
}

func (s *Shard) AddGuildMemberRole(ctx context.Context, guildId, userId, roleId uint64) error {
	return rest.AddGuildMemberRole(ctx, s.Token, s.ShardManager.RateLimiter, guildId, userId, roleId)
}

func (s *Shard) RemoveGuildMemberRole(ctx context.Context, guildId, userId, roleId uint64) error {
	return rest.RemoveGuildMemberRole(ctx, s.Token, s.ShardManager.RateLimiter, guildId, userId, roleId)
}

func (s *Shard) RemoveGuildMember(ctx context.Context, guildId, userId uint64) error {
	return rest.RemoveGuildMember(ctx, s.Token, s.ShardManager.RateLimiter, guildId, userId)
}

func (s *Shard) GetGuildBans(ctx context.Context, guildId uint64, data rest.GetGuildBansData) ([]guild.Ban, error) {
	return rest.GetGuildBans(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) GetGuildBan(ctx context.Context, guildId, userId uint64) (guild.Ban, error) {
	return rest.GetGuildBan(ctx, s.Token, s.ShardManager.RateLimiter, guildId, userId)
}

func (s *Shard) CreateGuildBan(ctx context.Context, guildId, userId uint64, data rest.CreateGuildBanData) error {
	return rest.CreateGuildBan(ctx, s.Token, s.ShardManager.RateLimiter, guildId, userId, data)
}

func (s *Shard) RemoveGuildBan(ctx context.Context, guildId, userId uint64) error {
	return rest.RemoveGuildBan(ctx, s.Token, s.ShardManager.RateLimiter, guildId, userId)
}

func (s *Shard) GetGuildRoles(ctx context.Context, guildId uint64) ([]guild.Role, error) {
	if s.Cache.Options().Roles {
		if cached, err := s.Cache.GetGuildRoles(ctx, guildId); err == nil {
			return cached, nil
		} else if err != cache.ErrNotFound {
			return nil, err
		}
	}

	roles, err := rest.GetGuildRoles(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
	if err != nil {
		return nil, err
	}

	if s.Cache.Options().Roles {
		if err := s.Cache.StoreRoles(ctx, roles, guildId); err != nil {
			return nil, err
		}
	}

	return roles, err
}

func (s *Shard) CreateGuildRole(ctx context.Context, guildId uint64, data rest.GuildRoleData) (guild.Role, error) {
	return rest.CreateGuildRole(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) ModifyGuildRolePositions(ctx context.Context, guildId uint64, positions []rest.Position) ([]guild.Role, error) {
	return rest.ModifyGuildRolePositions(ctx, s.Token, s.ShardManager.RateLimiter, guildId, positions)
}

func (s *Shard) ModifyGuildRole(ctx context.Context, guildId, roleId uint64, data rest.GuildRoleData) (guild.Role, error) {
	return rest.ModifyGuildRole(ctx, s.Token, s.ShardManager.RateLimiter, guildId, roleId, data)
}

func (s *Shard) DeleteGuildRole(ctx context.Context, guildId, roleId uint64) error {
	return rest.DeleteGuildRole(ctx, s.Token, s.ShardManager.RateLimiter, guildId, roleId)
}

func (s *Shard) GetGuildPruneCount(ctx context.Context, guildId uint64, days int) (int, error) {
	return rest.GetGuildPruneCount(ctx, s.Token, s.ShardManager.RateLimiter, guildId, days)
}

// computePruneCount = whether 'pruned' is returned, discouraged for large guilds
func (s *Shard) BeginGuildPrune(ctx context.Context, guildId uint64, days int, computePruneCount bool) error {
	return rest.BeginGuildPrune(ctx, s.Token, s.ShardManager.RateLimiter, guildId, days, computePruneCount)
}

func (s *Shard) GetGuildVoiceRegions(ctx context.Context, guildId uint64) ([]guild.VoiceRegion, error) {
	return rest.GetGuildVoiceRegions(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetGuildInvites(ctx context.Context, guildId uint64) ([]invite.InviteMetadata, error) {
	return rest.GetGuildInvites(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetGuildIntegrations(ctx context.Context, guildId uint64) ([]integration.Integration, error) {
	return rest.GetGuildIntegrations(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) CreateGuildIntegration(ctx context.Context, guildId uint64, data rest.CreateIntegrationData) error {
	return rest.CreateGuildIntegration(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) ModifyGuildIntegration(ctx context.Context, guildId, integrationId uint64, data rest.ModifyIntegrationData) error {
	return rest.ModifyGuildIntegration(ctx, s.Token, s.ShardManager.RateLimiter, guildId, integrationId, data)
}

func (s *Shard) DeleteGuildIntegration(ctx context.Context, guildId, integrationId uint64) error {
	return rest.DeleteGuildIntegration(ctx, s.Token, s.ShardManager.RateLimiter, guildId, integrationId)
}

func (s *Shard) SyncGuildIntegration(ctx context.Context, guildId, integrationId uint64) error {
	return rest.SyncGuildIntegration(ctx, s.Token, s.ShardManager.RateLimiter, guildId, integrationId)
}

func (s *Shard) GetGuildWidget(ctx context.Context, guildId uint64) (guild.GuildWidget, error) {
	return rest.GetGuildWidget(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) ModifyGuildEmbed(ctx context.Context, guildId uint64, data guild.GuildEmbed) (guild.GuildEmbed, error) {
	return rest.ModifyGuildEmbed(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
}

// returns invite object with only "code" and "uses" fields
func (s *Shard) GetGuildVanityUrl(ctx context.Context, guildId uint64) (invite.Invite, error) {
	return rest.GetGuildVanityURL(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetInvite(ctx context.Context, inviteCode string, withCounts bool) (invite.Invite, error) {
	return rest.GetInvite(ctx, s.Token, s.ShardManager.RateLimiter, inviteCode, withCounts)
}

func (s *Shard) DeleteInvite(ctx context.Context, inviteCode string) (invite.Invite, error) {
	return rest.DeleteInvite(ctx, s.Token, s.ShardManager.RateLimiter, inviteCode)
}

func (s *Shard) GetCurrentUser(ctx context.Context) (user.User, error) {
	if cached, err := s.Cache.GetSelf(ctx); err == nil {
		return cached, nil
	} else if err != cache.ErrNotFound {
		return user.User{}, err
	}

	self, err := rest.GetCurrentUser(ctx, s.Token, s.ShardManager.RateLimiter)
	if err != nil {
		return user.User{}, err
	}

	if err := s.Cache.StoreSelf(ctx, self); err != nil {
		return user.User{}, err
	}

	return self, err
}

func (s *Shard) GetUser(ctx context.Context, userId uint64) (user.User, error) {
	if s.Cache.Options().Users {
		if cached, err := s.Cache.GetUser(ctx, userId); err == nil {
			return cached, nil
		} else if err != cache.ErrNotFound {
			return user.User{}, err
		}
	}

	u, err := rest.GetUser(ctx, s.Token, s.ShardManager.RateLimiter, userId)
	if err != nil {
		return user.User{}, err
	}

	if s.Cache.Options().Users {
		if err := s.Cache.StoreUser(ctx, u); err != nil {
			return user.User{}, err
		}
	}

	return u, err
}

func (s *Shard) ModifyCurrentUser(ctx context.Context, data rest.ModifyUserData) (user.User, error) {
	return rest.ModifyCurrentUser(ctx, s.Token, s.ShardManager.RateLimiter, data)
}

func (s *Shard) GetCurrentUserGuilds(ctx context.Context, data rest.CurrentUserGuildsData) ([]guild.Guild, error) {
	return rest.GetCurrentUserGuilds(ctx, s.Token, s.ShardManager.RateLimiter, data)
}

func (s *Shard) LeaveGuild(ctx context.Context, guildId uint64) error {
	return rest.LeaveGuild(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) CreateDM(ctx context.Context, recipientId uint64) (channel.Channel, error) {
	return rest.CreateDM(ctx, s.Token, s.ShardManager.RateLimiter, recipientId)
}

func (s *Shard) GetUserConnections(ctx context.Context) ([]integration.Connection, error) {
	return rest.GetUserConnections(ctx, s.Token, s.ShardManager.RateLimiter)
}

// GetGuildVoiceRegions should be preferred, as it returns VIP servers if available to the guild
func (s *Shard) ListVoiceRegions(ctx context.Context) ([]guild.VoiceRegion, error) {
	return rest.ListVoiceRegions(ctx, s.Token)
}

func (s *Shard) CreateWebhook(ctx context.Context, channelId uint64, data rest.WebhookData) (guild.Webhook, error) {
	return rest.CreateWebhook(ctx, s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) GetChannelWebhooks(ctx context.Context, channelId uint64) ([]guild.Webhook, error) {
	return rest.GetChannelWebhooks(ctx, s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) GetGuildWebhooks(ctx context.Context, guildId uint64) ([]guild.Webhook, error) {
	return rest.GetGuildWebhooks(ctx, s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetWebhook(ctx context.Context, webhookId uint64) (guild.Webhook, error) {
	return rest.GetWebhook(ctx, s.Token, s.ShardManager.RateLimiter, webhookId)
}

func (s *Shard) ModifyWebhook(ctx context.Context, webhookId uint64, data rest.ModifyWebhookData) (guild.Webhook, error) {
	return rest.ModifyWebhook(ctx, s.Token, s.ShardManager.RateLimiter, webhookId, data)
}

func (s *Shard) DeleteWebhook(ctx context.Context, webhookId uint64) error {
	return rest.DeleteWebhook(ctx, s.Token, s.ShardManager.RateLimiter, webhookId)
}

// if wait=true, a message object will be returned
func (s *Shard) ExecuteWebhook(ctx context.Context, webhookId uint64, webhookToken string, wait bool, data rest.WebhookBody) (*message.Message, error) {
	return rest.ExecuteWebhook(ctx, webhookToken, s.ShardManager.RateLimiter, webhookId, wait, data)
}

func (s *Shard) EditWebhookMessage(ctx context.Context, webhookId uint64, webhookToken string, messageId uint64, data rest.WebhookEditBody) (message.Message, error) {
	return rest.EditWebhookMessage(ctx, webhookToken, s.ShardManager.RateLimiter, webhookId, messageId, data)
}

func (s *Shard) GetGuildAuditLog(ctx context.Context, guildId uint64, data rest.GetGuildAuditLogData) (auditlog.AuditLog, error) {
	return rest.GetGuildAuditLog(ctx, s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) GetGlobalCommands(ctx context.Context, applicationId uint64) ([]interaction.ApplicationCommand, error) {
	return rest.GetGlobalCommands(ctx, s.Token, s.ShardManager.RateLimiter, applicationId)
}

func (s *Shard) CreateGlobalCommand(ctx context.Context, applicationId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.CreateGlobalCommand(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, data)
}

func (s *Shard) ModifyGlobalCommand(ctx context.Context, applicationId, commandId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.ModifyGlobalCommand(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, commandId, data)
}

func (s *Shard) ModifyGlobalCommands(ctx context.Context, applicationId uint64, data []rest.CreateCommandData) ([]interaction.ApplicationCommand, error) {
	return rest.ModifyGlobalCommands(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, data)
}

func (s *Shard) DeleteGlobalCommand(ctx context.Context, applicationId, commandId uint64) error {
	return rest.DeleteGlobalCommand(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, commandId)
}

func (s *Shard) GetGuildCommands(ctx context.Context, applicationId, guildId uint64) ([]interaction.ApplicationCommand, error) {
	return rest.GetGuildCommands(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId)
}

func (s *Shard) CreateGuildCommand(ctx context.Context, applicationId, guildId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.CreateGuildCommand(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId, data)
}

func (s *Shard) ModifyGuildCommand(ctx context.Context, applicationId, guildId, commandId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.ModifyGuildCommand(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId, commandId, data)
}

func (s *Shard) ModifyGuildCommands(ctx context.Context, applicationId, guildId uint64, data []rest.CreateCommandData) ([]interaction.ApplicationCommand, error) {
	return rest.ModifyGuildCommands(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId, data)
}

func (s *Shard) DeleteGuildCommand(ctx context.Context, applicationId, guildId, commandId uint64) error {
	return rest.DeleteGuildCommand(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId, commandId)
}

func (s *Shard) GetCommandPermissions(ctx context.Context, applicationId, guildId, commandId uint64) (rest.CommandWithPermissionsData, error) {
	return rest.GetCommandPermissions(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId, commandId)
}

func (s *Shard) GetBulkCommandPermissions(ctx context.Context, applicationId, guildId uint64) ([]rest.CommandWithPermissionsData, error) {
	return rest.GetBulkCommandPermissions(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId)
}

func (s *Shard) EditCommandPermissions(ctx context.Context, applicationId, guildId, commandId uint64, data rest.CommandWithPermissionsData) (rest.CommandWithPermissionsData, error) {
	return rest.EditCommandPermissions(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId, commandId, data)
}

func (s *Shard) EditBulkCommandPermissions(ctx context.Context, applicationId, guildId uint64, data []rest.CommandWithPermissionsData) ([]rest.CommandWithPermissionsData, error) {
	return rest.EditBulkCommandPermissions(ctx, s.Token, s.ShardManager.RateLimiter, applicationId, guildId, data)
}
