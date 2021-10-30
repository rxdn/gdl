package gateway

import (
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

func (s *Shard) GetChannel(channelId uint64) (channel.Channel, error) {
	shouldCache := s.Cache.GetOptions().Channels
	if shouldCache {
		if cached, found := s.Cache.GetChannel(channelId); found {
			return cached, nil
		}
	}

	channel, err := rest.GetChannel(s.Token, s.ShardManager.RateLimiter, channelId)

	if shouldCache && err == nil {
		go s.Cache.StoreChannel(channel)
	}

	return channel, err
}

func (s *Shard) ModifyChannel(channelId uint64, data rest.ModifyChannelData) (channel.Channel, error) {
	channel, err := rest.ModifyChannel(s.Token, s.ShardManager.RateLimiter, channelId, data)

	if s.Cache.GetOptions().Channels && err != nil {
		go s.Cache.StoreChannel(channel)
	}

	return channel, err
}

func (s *Shard) DeleteChannel(channelId uint64) (channel.Channel, error) {
	return rest.DeleteChannel(s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) GetChannelMessages(channelId uint64, options rest.GetChannelMessagesData) ([]message.Message, error) {
	return rest.GetChannelMessages(s.Token, s.ShardManager.RateLimiter, channelId, options)
}

func (s *Shard) GetChannelMessage(channelId, messageId uint64) (message.Message, error) {
	return rest.GetChannelMessage(s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) CreateMessage(channelId uint64, content string) (message.Message, error) {
	return s.CreateMessageComplex(channelId, rest.CreateMessageData{
		Content: content,
	})
}

func (s *Shard) CreateMessageReply(channelId uint64, content string, reference *message.MessageReference) (message.Message, error) {
	return s.CreateMessageComplex(channelId, rest.CreateMessageData{
		Content:          content,
		MessageReference: reference,
	})
}

func (s *Shard) CreateMessageEmbed(channelId uint64, embed ...*embed.Embed) (message.Message, error) {
	return s.CreateMessageComplex(channelId, rest.CreateMessageData{
		Embeds: embed,
	})
}

func (s *Shard) CreateMessageEmbedReply(channelId uint64, e *embed.Embed, reference *message.MessageReference) (message.Message, error) {
	return s.CreateMessageComplex(channelId, rest.CreateMessageData{
		Embeds:           []*embed.Embed{e},
		MessageReference: reference,
	})
}

func (s *Shard) CreateMessageComplex(channelId uint64, data rest.CreateMessageData) (message.Message, error) {
	return rest.CreateMessage(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) CreateReaction(channelId, messageId uint64, emoji string) error {
	return rest.CreateReaction(s.Token, s.ShardManager.RateLimiter, channelId, messageId, emoji)
}

func (s *Shard) DeleteOwnReaction(channelId, messageId uint64, emoji string) error {
	return rest.DeleteOwnReaction(s.Token, s.ShardManager.RateLimiter, channelId, messageId, emoji)
}

func (s *Shard) DeleteUserReaction(channelId, messageId, userId uint64, emoji string) error {
	return rest.DeleteUserReaction(s.Token, s.ShardManager.RateLimiter, channelId, messageId, userId, emoji)
}

func (s *Shard) GetReactions(channelId, messageId uint64, emoji string, options rest.GetReactionsData) ([]user.User, error) {
	return rest.GetReactions(s.Token, s.ShardManager.RateLimiter, channelId, messageId, emoji, options)
}

func (s *Shard) DeleteAllReactions(channelId, messageId uint64) error {
	return rest.DeleteAllReactions(s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) DeleteAllReactionsEmoji(channelId, messageId uint64, emoji string) error {
	return rest.DeleteAllReactionsEmoji(s.Token, s.ShardManager.RateLimiter, channelId, messageId, emoji)
}

func (s *Shard) EditMessage(channelId, messageId uint64, data rest.EditMessageData) (message.Message, error) {
	return rest.EditMessage(s.Token, s.ShardManager.RateLimiter, channelId, messageId, data)
}

func (s *Shard) DeleteMessage(channelId, messageId uint64) error {
	return rest.DeleteMessage(s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) BulkDeleteMessages(channelId uint64, messages []uint64) error {
	return rest.BulkDeleteMessages(s.Token, s.ShardManager.RateLimiter, channelId, messages)
}

func (s *Shard) EditChannelPermissions(channelId uint64, updated channel.PermissionOverwrite) error {
	return rest.EditChannelPermissions(s.Token, s.ShardManager.RateLimiter, channelId, updated)
}

func (s *Shard) GetChannelInvites(channelId uint64) ([]invite.InviteMetadata, error) {
	return rest.GetChannelInvites(s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) CreateChannelInvite(channelId uint64, data rest.CreateInviteData) (invite.Invite, error) {
	return rest.CreateChannelInvite(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) DeleteChannelPermissions(channelId, overwriteId uint64) error {
	return rest.DeleteChannelPermissions(s.Token, s.ShardManager.RateLimiter, channelId, overwriteId)
}

func (s *Shard) TriggerTypingIndicator(channelId uint64) error {
	return rest.TriggerTypingIndicator(s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) GetPinnedMessages(channelId uint64) ([]message.Message, error) {
	return rest.GetPinnedMessages(s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) AddPinnedChannelMessage(channelId, messageId uint64) error {
	return rest.AddPinnedChannelMessage(s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) DeletePinnedChannelMessage(channelId, messageId uint64) error {
	return rest.DeletePinnedChannelMessage(s.Token, s.ShardManager.RateLimiter, channelId, messageId)
}

func (s *Shard) ListThreadMembers(channelId uint64) ([]channel.ThreadMember, error) {
	return rest.ListThreadMembers(s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) ListActiveThreads(channelId uint64) (rest.ThreadsResponse, error) {
	return rest.ListActiveThreads(s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) ListPublicArchivedThreads(channelId uint64, data rest.ListThreadsData) (rest.ThreadsResponse, error) {
	return rest.ListPublicArchivedThreads(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) ListPrivateArchivedThreads(channelId uint64, data rest.ListThreadsData) (rest.ThreadsResponse, error) {
	return rest.ListPrivateArchivedThreads(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) ListJoinedPrivateArchivedThreads(channelId uint64, data rest.ListThreadsData) (rest.ThreadsResponse, error) {
	return rest.ListPrivateArchivedThreads(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) StartThreadWithMessage(channelId, messageId uint64, data rest.StartThreadWithMessageData) (channel.Channel, error) {
	return rest.StartThreadWithMessage(s.Token, s.ShardManager.RateLimiter, channelId, messageId, data)
}

func (s *Shard) StartThreadWithoutMessage(channelId, messageId uint64, data rest.StartThreadWithoutMessageData) (channel.Channel, error) {
	return rest.StartThreadWithoutMessage(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) CreatePublicThread(channelId uint64, name string, autoArchiveDuration uint16) (channel.Channel, error) {
	data := rest.StartThreadWithoutMessageData{
		Name:                name,
		AutoArchiveDuration: autoArchiveDuration,
		Type:                channel.ChannelTypeGuildPublicThread,
	}

	return rest.StartThreadWithoutMessage(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) CreatePrivateThread(channelId uint64, name string, autoArchiveDuration uint16, invitable bool) (channel.Channel, error) {
	data := rest.StartThreadWithoutMessageData{
		Name:                name,
		AutoArchiveDuration: autoArchiveDuration,
		Type:                channel.ChannelTypeGuildPrivateThread,
		Invitable:           invitable,
	}

	return rest.StartThreadWithoutMessage(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) ListGuildEmojis(guildId uint64) ([]emoji.Emoji, error) {
	shouldCacheEmoji := s.Cache.GetOptions().Emojis
	shouldCacheGuild := s.Cache.GetOptions().Guilds

	if shouldCacheEmoji && shouldCacheGuild {
		if guild, found := s.Cache.GetGuild(guildId, false); found {
			return guild.Emojis, nil
		}
	}

	emojis, err := rest.ListGuildEmojis(s.Token, s.ShardManager.RateLimiter, guildId)

	if shouldCacheEmoji && err == nil {
		go s.Cache.StoreEmojis(emojis, guildId)
	}

	return emojis, err
}

func (s *Shard) GetGuildEmoji(guildId uint64, emojiId uint64) (emoji.Emoji, error) {
	shouldCache := s.Cache.GetOptions().Emojis
	if shouldCache {
		if emoji, found := s.Cache.GetEmoji(emojiId); found {
			return emoji, nil
		}
	}

	emoji, err := rest.GetGuildEmoji(s.Token, s.ShardManager.RateLimiter, guildId, emojiId)

	if shouldCache && err == nil {
		go s.Cache.StoreEmoji(emoji, guildId)
	}

	return emoji, err
}

func (s *Shard) CreateGuildEmoji(guildId uint64, data rest.CreateEmojiData) (emoji.Emoji, error) {
	return rest.CreateGuildEmoji(s.Token, s.ShardManager.RateLimiter, guildId, data)
}

// updating Image is not permitted
func (s *Shard) ModifyGuildEmoji(guildId, emojiId uint64, data rest.CreateEmojiData) (emoji.Emoji, error) {
	return rest.ModifyGuildEmoji(s.Token, s.ShardManager.RateLimiter, guildId, emojiId, data)
}

func (s *Shard) CreateGuild(data rest.CreateGuildData) (guild.Guild, error) {
	return rest.CreateGuild(s.Token, data)
}

func (s *Shard) GetGuild(guildId uint64) (guild.Guild, error) {
	shouldCache := s.Cache.GetOptions().Guilds

	if shouldCache {
		if cachedGuild, found := s.Cache.GetGuild(guildId, false); found {
			return cachedGuild, nil
		}
	}

	guild, err := rest.GetGuild(s.Token, s.ShardManager.RateLimiter, guildId)
	if err == nil {
		go s.Cache.StoreGuild(guild)
	}

	return guild, err
}

func (s *Shard) GetGuildPreview(guildId uint64) (guild.GuildPreview, error) {
	return rest.GetGuildPreview(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) ModifyGuild(guildId uint64, data rest.ModifyGuildData) (guild.Guild, error) {
	return rest.ModifyGuild(s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) DeleteGuild(guildId uint64) error {
	return rest.DeleteGuild(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetGuildChannels(guildId uint64) ([]channel.Channel, error) {
	shouldCache := s.Cache.GetOptions().Guilds && s.Cache.GetOptions().Channels

	if shouldCache {
		cached := s.Cache.GetGuildChannels(guildId)

		// either not cached (more likely), or guild has no channels
		if len(cached) > 0 {
			return cached, nil
		}
	}

	channels, err := rest.GetGuildChannels(s.Token, s.ShardManager.RateLimiter, guildId)

	if shouldCache && err == nil {
		go s.Cache.StoreChannels(channels)
	}

	return channels, err
}

func (s *Shard) CreateGuildChannel(guildId uint64, data rest.CreateChannelData) (channel.Channel, error) {
	return rest.CreateGuildChannel(s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) ModifyGuildChannelPositions(guildId uint64, positions []rest.Position) error {
	return rest.ModifyGuildChannelPositions(s.Token, s.ShardManager.RateLimiter, guildId, positions)
}

func (s *Shard) GetGuildMember(guildId, userId uint64) (member.Member, error) {
	cacheGuilds := s.Cache.GetOptions().Guilds
	cacheUsers := s.Cache.GetOptions().Users

	if cacheGuilds && cacheUsers {
		if member, found := s.Cache.GetMember(guildId, userId); found {
			return member, nil
		}
	}

	member, err := rest.GetGuildMember(s.Token, s.ShardManager.RateLimiter, guildId, userId)

	if cacheGuilds && err == nil {
		go s.Cache.StoreMember(member, guildId)
	}

	return member, err
}

func (s *Shard) SearchGuildMembers(guildId uint64, data rest.SearchGuildMembersData) ([]member.Member, error) {
	members, err := rest.SearchGuildMembers(s.Token, s.ShardManager.RateLimiter, guildId, data)
	if err == nil {
		go s.Cache.StoreMembers(members, guildId)
	}

	return members, err
}

func (s *Shard) ListGuildMembers(guildId uint64, data rest.ListGuildMembersData) ([]member.Member, error) {
	members, err := rest.ListGuildMembers(s.Token, s.ShardManager.RateLimiter, guildId, data)
	if err == nil {
		go s.Cache.StoreMembers(members, guildId)
	}

	return members, err
}

func (s *Shard) ModifyGuildMember(guildId, userId uint64, data rest.ModifyGuildMemberData) error {
	return rest.ModifyGuildMember(s.Token, s.ShardManager.RateLimiter, guildId, userId, data)
}

func (s *Shard) ModifyCurrentUserNick(guildId uint64, nick string) error {
	return rest.ModifyCurrentUserNick(s.Token, s.ShardManager.RateLimiter, guildId, nick)
}

func (s *Shard) AddGuildMemberRole(guildId, userId, roleId uint64) error {
	return rest.AddGuildMemberRole(s.Token, s.ShardManager.RateLimiter, guildId, userId, roleId)
}

func (s *Shard) RemoveGuildMemberRole(guildId, userId, roleId uint64) error {
	return rest.RemoveGuildMemberRole(s.Token, s.ShardManager.RateLimiter, guildId, userId, roleId)
}

func (s *Shard) RemoveGuildMember(guildId, userId uint64) error {
	return rest.RemoveGuildMember(s.Token, s.ShardManager.RateLimiter, guildId, userId)
}

func (s *Shard) GetGuildBans(guildId uint64) ([]guild.Ban, error) {
	return rest.GetGuildBans(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetGuildBan(guildId, userId uint64) (guild.Ban, error) {
	return rest.GetGuildBan(s.Token, s.ShardManager.RateLimiter, guildId, userId)
}

func (s *Shard) CreateGuildBan(guildId, userId uint64, data rest.CreateGuildBanData) error {
	return rest.CreateGuildBan(s.Token, s.ShardManager.RateLimiter, guildId, userId, data)
}

func (s *Shard) RemoveGuildBan(guildId, userId uint64) error {
	return rest.RemoveGuildBan(s.Token, s.ShardManager.RateLimiter, guildId, userId)
}

func (s *Shard) GetGuildRoles(guildId uint64) ([]guild.Role, error) {
	shouldCache := s.Cache.GetOptions().Guilds
	if shouldCache {
		cached := s.Cache.GetGuildRoles(guildId)

		// either not cached (more likely), or guild has no channels
		if len(cached) > 0 {
			return cached, nil
		}
	}

	roles, err := rest.GetGuildRoles(s.Token, s.ShardManager.RateLimiter, guildId)

	if shouldCache && err == nil {
		go s.Cache.StoreRoles(roles, guildId)
	}

	return roles, err
}

func (s *Shard) CreateGuildRole(guildId uint64, data rest.GuildRoleData) (guild.Role, error) {
	return rest.CreateGuildRole(s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) ModifyGuildRolePositions(guildId uint64, positions []rest.Position) ([]guild.Role, error) {
	return rest.ModifyGuildRolePositions(s.Token, s.ShardManager.RateLimiter, guildId, positions)
}

func (s *Shard) ModifyGuildRole(guildId, roleId uint64, data rest.GuildRoleData) (guild.Role, error) {
	return rest.ModifyGuildRole(s.Token, s.ShardManager.RateLimiter, guildId, roleId, data)
}

func (s *Shard) DeleteGuildRole(guildId, roleId uint64) error {
	return rest.DeleteGuildRole(s.Token, s.ShardManager.RateLimiter, guildId, roleId)
}

func (s *Shard) GetGuildPruneCount(guildId uint64, days int) (int, error) {
	return rest.GetGuildPruneCount(s.Token, s.ShardManager.RateLimiter, guildId, days)
}

// computePruneCount = whether 'pruned' is returned, discouraged for large guilds
func (s *Shard) BeginGuildPrune(guildId uint64, days int, computePruneCount bool) error {
	return rest.BeginGuildPrune(s.Token, s.ShardManager.RateLimiter, guildId, days, computePruneCount)
}

func (s *Shard) GetGuildVoiceRegions(guildId uint64) ([]guild.VoiceRegion, error) {
	return rest.GetGuildVoiceRegions(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetGuildInvites(guildId uint64) ([]invite.InviteMetadata, error) {
	return rest.GetGuildInvites(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetGuildIntegrations(guildId uint64) ([]integration.Integration, error) {
	return rest.GetGuildIntegrations(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) CreateGuildIntegration(guildId uint64, data rest.CreateIntegrationData) error {
	return rest.CreateGuildIntegration(s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) ModifyGuildIntegration(guildId, integrationId uint64, data rest.ModifyIntegrationData) error {
	return rest.ModifyGuildIntegration(s.Token, s.ShardManager.RateLimiter, guildId, integrationId, data)
}

func (s *Shard) DeleteGuildIntegration(guildId, integrationId uint64) error {
	return rest.DeleteGuildIntegration(s.Token, s.ShardManager.RateLimiter, guildId, integrationId)
}

func (s *Shard) SyncGuildIntegration(guildId, integrationId uint64) error {
	return rest.SyncGuildIntegration(s.Token, s.ShardManager.RateLimiter, guildId, integrationId)
}

func (s *Shard) GetGuildWidget(guildId uint64) (guild.GuildWidget, error) {
	return rest.GetGuildWidget(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) ModifyGuildEmbed(guildId uint64, data guild.GuildEmbed) (guild.GuildEmbed, error) {
	return rest.ModifyGuildEmbed(s.Token, s.ShardManager.RateLimiter, guildId, data)
}

// returns invite object with only "code" and "uses" fields
func (s *Shard) GetGuildVanityUrl(guildId uint64) (invite.Invite, error) {
	return rest.GetGuildVanityURL(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetInvite(inviteCode string, withCounts bool) (invite.Invite, error) {
	return rest.GetInvite(s.Token, s.ShardManager.RateLimiter, inviteCode, withCounts)
}

func (s *Shard) DeleteInvite(inviteCode string) (invite.Invite, error) {
	return rest.DeleteInvite(s.Token, s.ShardManager.RateLimiter, inviteCode)
}

func (s *Shard) GetCurrentUser() (user.User, error) {
	if cached, found := s.Cache.GetSelf(); found {
		return cached, nil
	}

	self, err := rest.GetCurrentUser(s.Token, s.ShardManager.RateLimiter)

	if err == nil {
		go s.Cache.StoreSelf(self)
	}

	return self, err
}

func (s *Shard) GetUser(userId uint64) (user.User, error) {
	shouldCache := s.Cache.GetOptions().Users

	if shouldCache {
		if cached, found := s.Cache.GetUser(userId); found {
			return cached, nil
		}
	}

	user, err := rest.GetUser(s.Token, s.ShardManager.RateLimiter, userId)

	if shouldCache && err == nil {
		go s.Cache.StoreUser(user)
	}

	return user, err
}

func (s *Shard) ModifyCurrentUser(data rest.ModifyUserData) (user.User, error) {
	return rest.ModifyCurrentUser(s.Token, s.ShardManager.RateLimiter, data)
}

func (s *Shard) GetCurrentUserGuilds(data rest.CurrentUserGuildsData) ([]guild.Guild, error) {
	return rest.GetCurrentUserGuilds(s.Token, s.ShardManager.RateLimiter, data)
}

func (s *Shard) LeaveGuild(guildId uint64) error {
	return rest.LeaveGuild(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) CreateDM(recipientId uint64) (channel.Channel, error) {
	return rest.CreateDM(s.Token, s.ShardManager.RateLimiter, recipientId)
}

func (s *Shard) GetUserConnections() ([]integration.Connection, error) {
	return rest.GetUserConnections(s.Token, s.ShardManager.RateLimiter)
}

// GetGuildVoiceRegions should be preferred, as it returns VIP servers if available to the guild
func (s *Shard) ListVoiceRegions() ([]guild.VoiceRegion, error) {
	return rest.ListVoiceRegions(s.Token)
}

func (s *Shard) CreateWebhook(channelId uint64, data rest.WebhookData) (guild.Webhook, error) {
	return rest.CreateWebhook(s.Token, s.ShardManager.RateLimiter, channelId, data)
}

func (s *Shard) GetChannelWebhooks(channelId uint64) ([]guild.Webhook, error) {
	return rest.GetChannelWebhooks(s.Token, s.ShardManager.RateLimiter, channelId)
}

func (s *Shard) GetGuildWebhooks(guildId uint64) ([]guild.Webhook, error) {
	return rest.GetGuildWebhooks(s.Token, s.ShardManager.RateLimiter, guildId)
}

func (s *Shard) GetWebhook(webhookId uint64) (guild.Webhook, error) {
	return rest.GetWebhook(s.Token, s.ShardManager.RateLimiter, webhookId)
}

func (s *Shard) ModifyWebhook(webhookId uint64, data rest.ModifyWebhookData) (guild.Webhook, error) {
	return rest.ModifyWebhook(s.Token, s.ShardManager.RateLimiter, webhookId, data)
}

func (s *Shard) DeleteWebhook(webhookId uint64) error {
	return rest.DeleteWebhook(s.Token, s.ShardManager.RateLimiter, webhookId)
}

// if wait=true, a message object will be returned
func (s *Shard) ExecuteWebhook(webhookId uint64, webhookToken string, wait bool, data rest.WebhookBody) (*message.Message, error) {
	return rest.ExecuteWebhook(webhookToken, s.ShardManager.RateLimiter, webhookId, wait, data)
}

func (s *Shard) EditWebhookMessage(webhookId uint64, webhookToken string, messageId uint64, data rest.WebhookEditBody) (message.Message, error) {
	return rest.EditWebhookMessage(webhookToken, s.ShardManager.RateLimiter, webhookId, messageId, data)
}

func (s *Shard) GetGuildAuditLog(guildId uint64, data rest.GetGuildAuditLogData) (auditlog.AuditLog, error) {
	return rest.GetGuildAuditLog(s.Token, s.ShardManager.RateLimiter, guildId, data)
}

func (s *Shard) GetGlobalCommands(applicationId uint64) ([]interaction.ApplicationCommand, error) {
	return rest.GetGlobalCommands(s.Token, s.ShardManager.RateLimiter, applicationId)
}

func (s *Shard) CreateGlobalCommand(applicationId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.CreateGlobalCommand(s.Token, s.ShardManager.RateLimiter, applicationId, data)
}

func (s *Shard) ModifyGlobalCommand(applicationId, commandId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.ModifyGlobalCommand(s.Token, s.ShardManager.RateLimiter, applicationId, commandId, data)
}

func (s *Shard) ModifyGlobalCommands(applicationId uint64, data []rest.CreateCommandData) ([]interaction.ApplicationCommand, error) {
	return rest.ModifyGlobalCommands(s.Token, s.ShardManager.RateLimiter, applicationId, data)
}

func (s *Shard) DeleteGlobalCommand(applicationId, commandId uint64) error {
	return rest.DeleteGlobalCommand(s.Token, s.ShardManager.RateLimiter, applicationId, commandId)
}

func (s *Shard) GetGuildCommands(applicationId, guildId uint64) ([]interaction.ApplicationCommand, error) {
	return rest.GetGuildCommands(s.Token, s.ShardManager.RateLimiter, applicationId, guildId)
}

func (s *Shard) CreateGuildCommand(applicationId, guildId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.CreateGuildCommand(s.Token, s.ShardManager.RateLimiter, applicationId, guildId, data)
}

func (s *Shard) ModifyGuildCommand(applicationId, guildId, commandId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.ModifyGuildCommand(s.Token, s.ShardManager.RateLimiter, applicationId, guildId, commandId, data)
}

func (s *Shard) ModifyGuildCommands(applicationId, guildId uint64, data []rest.CreateCommandData) ([]interaction.ApplicationCommand, error) {
	return rest.ModifyGuildCommands(s.Token, s.ShardManager.RateLimiter, applicationId, guildId, data)
}

func (s *Shard) DeleteGuildCommand(applicationId, guildId, commandId uint64) error {
	return rest.DeleteGuildCommand(s.Token, s.ShardManager.RateLimiter, applicationId, guildId, commandId)
}

func (s *Shard) GetCommandPermissions(applicationId, guildId, commandId uint64) (rest.CommandWithPermissionsData, error) {
	return rest.GetCommandPermissions(s.Token, s.ShardManager.RateLimiter, applicationId, guildId, commandId)
}

func (s *Shard) GetBulkCommandPermissions(applicationId, guildId uint64) ([]rest.CommandWithPermissionsData, error) {
	return rest.GetBulkCommandPermissions(s.Token, s.ShardManager.RateLimiter, applicationId, guildId)
}

func (s *Shard) EditCommandPermissions(applicationId, guildId, commandId uint64, data rest.CommandWithPermissionsData) (rest.CommandWithPermissionsData, error) {
	return rest.EditCommandPermissions(s.Token, s.ShardManager.RateLimiter, applicationId, guildId, commandId, data)
}

func (s *Shard) EditBulkCommandPermissions(applicationId, guildId uint64, data []rest.CommandWithPermissionsData) ([]rest.CommandWithPermissionsData, error) {
	return rest.EditBulkCommandPermissions(s.Token, s.ShardManager.RateLimiter, applicationId, guildId, data)
}
