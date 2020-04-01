package gateway

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/integration"
	"github.com/rxdn/gdl/objects/invite"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest"
	"image"
)

func (s *Shard) GetChannel(channelId uint64) (*channel.Channel, error) {
	shouldCache := (*s.Cache).GetOptions().Channels
	if shouldCache {
		cached := (*s.Cache).GetChannel(channelId)
		if cached != nil {
			return cached, nil
		}
	}

	channel, err := rest.GetChannel(s.Token, channelId)

	if shouldCache && err == nil {
		go func() {
			lock := (*s.Cache).GetLock(channelId)
			lock.Lock()
			(*s.Cache).StoreChannel(channel)
			lock.Unlock()
		}()
	}

	return channel, err
}

func (s *Shard) ModifyChannel(channelId uint64, data rest.ModifyChannelData) (*channel.Channel, error) {
	channel, err := rest.ModifyChannel(s.Token, channelId, data)

	if (*s.Cache).GetOptions().Channels && err != nil {
		lock := (*s.Cache).GetLock(channelId)
		lock.Lock()
		(*s.Cache).StoreChannel(channel)
		lock.Unlock()
	}

	return channel, err
}

func (s *Shard) DeleteChannel(channelId uint64) (*channel.Channel, error) {
	return rest.DeleteChannel(s.Token, channelId)
}

func (s *Shard) GetChannelMessages(channelId uint64, options rest.GetChannelMessagesData) ([]*message.Message, error) {
	return rest.GetChannelMessages(s.Token, channelId, options)
}

func (s *Shard) GetChannelMessage(channelId, messageId uint64) (*message.Message, error) {
	return rest.GetChannelMessage(s.Token, channelId, messageId)
}

func (s *Shard) CreateMessage(channelId uint64, content string) (*message.Message, error) {
	return s.CreateMessageComplex(channelId, rest.CreateMessageData{
		Content: content,
	})
}

func (s *Shard) CreateMessageEmbed(channelId uint64, embed *embed.Embed) (*message.Message, error) {
	return s.CreateMessageComplex(channelId, rest.CreateMessageData{
		Embed: embed,
	})
}

func (s *Shard) CreateMessageComplex(channelId uint64, data rest.CreateMessageData) (*message.Message, error) {
	return rest.CreateMessage(s.Token, channelId, data)
}

func (s *Shard) CreateReaction(channelId, messageId uint64, emoji string) error {
	return rest.CreateReaction(s.Token, channelId, messageId, emoji)
}

func (s *Shard) DeleteOwnReaction(channelId, messageId uint64, emoji string) error {
	return rest.DeleteOwnReaction(s.Token, channelId, messageId, emoji)
}

func (s *Shard) DeleteUserReaction(channelId, messageId, userId uint64, emoji string) error {
	return rest.DeleteUserReaction(s.Token, channelId, messageId, userId, emoji)
}

func (s *Shard) GetReactions(channelId, messageId uint64, emoji string, options rest.GetReactionsData) ([]user.User, error) {
	return rest.GetReactions(s.Token, channelId, messageId, emoji, options)
}

func (s *Shard) DeleteAllReactions(channelId, messageId uint64) error {
	return rest.DeleteAllReactions(s.Token, channelId, messageId)
}

func (s *Shard) DeleteAllReactionsEmoji(channelId, messageId uint64, emoji string) error {
	return rest.DeleteAllReactionsEmoji(s.Token, channelId, messageId, emoji)
}

func (s *Shard) EditMessage(channelId, messageId uint64, data rest.ModifyChannelData) (*message.Message, error) {
	return rest.EditMessage(s.Token, channelId, messageId, data)
}

func (s *Shard) DeleteMessage(channelId, messageId uint64) error {
	return rest.DeleteMessage(s.Token, channelId, messageId)
}

func (s *Shard) BulkDeleteMessages(channelId uint64, messages []uint64) error {
	return rest.BulkDeleteMessages(s.Token, channelId, messages)
}

func (s *Shard) EditChannelPermissions(channelId uint64, updated channel.PermissionOverwrite) error {
	return rest.EditChannelPermissions(s.Token, channelId, updated)
}

func (s *Shard) GetChannelInvites(channelId uint64) ([]invite.InviteMetadata, error) {
	return rest.GetChannelInvites(s.Token, channelId)
}

func (s *Shard) CreateChannelInvite(channelId uint64, data invite.InviteMetadata) (*invite.Invite, error) {
	return rest.CreateChannelInvite(s.Token, channelId, data)
}

func (s *Shard) DeleteChannelPermissions(channelId, overwriteId uint64) error {
	return rest.DeleteChannelPermissions(s.Token, channelId, overwriteId)
}

func (s *Shard) TriggerTypingIndicator(channelId uint64) error {
	return rest.TriggerTypingIndicator(s.Token, channelId)
}

func (s *Shard) GetPinnedMessages(channelId uint64) ([]*message.Message, error) {
	return rest.GetPinnedMessages(s.Token, channelId)
}

func (s *Shard) AddPinnedChannelMessage(channelId, messageId uint64) error {
	return rest.AddPinnedChannelMessage(s.Token, channelId, messageId)
}

func (s *Shard) DeletePinnedChannelMessage(channelId, messageId uint64) error {
	return rest.DeletePinnedChannelMessage(s.Token, channelId, messageId)
}

func (s *Shard) ListGuildEmojis(guildId uint64) ([]*emoji.Emoji, error) {
	shouldCacheEmoji := (*s.Cache).GetOptions().Emojis
	shouldCacheGuild := (*s.Cache).GetOptions().Guilds

	if shouldCacheEmoji && shouldCacheGuild {
		guild := (*s.Cache).GetGuild(guildId)
		if guild != nil {
			return guild.Emojis, nil
		}
	}

	emojis, err := rest.ListGuildEmojis(s.Token, guildId)

	if shouldCacheEmoji && err == nil {
		go func() {
			for _, emoji := range emojis {
				(*s.Cache).StoreEmoji(emoji)
			}

			if shouldCacheGuild {
				lock := (*s.Cache).GetLock(guildId)
				lock.Lock()

				guild := (*s.Cache).GetGuild(guildId)
				if guild != nil {
					guild.Emojis = emojis
					(*s.Cache).StoreGuild(guild)
				}

				lock.Unlock()
			}
		}()
	}

	return emojis, err
}

func (s *Shard) GetGuildEmoji(guildId uint64, emojiId uint64) (*emoji.Emoji, error) {
	shouldCache := (*s.Cache).GetOptions().Emojis
	if shouldCache {
		emoji := (*s.Cache).GetEmoji(emojiId)
		if emoji != nil {
			return emoji, nil
		}
	}

	emoji, err := rest.GetGuildEmoji(s.Token, guildId, emojiId)

	if shouldCache && err == nil {
		lock := (*s.Cache).GetLock(emojiId)
		lock.Lock()
		(*s.Cache).StoreEmoji(emoji)
		lock.Unlock()
	}

	return emoji, err
}

func (s *Shard) CreateGuildEmoji(guildId uint64, data rest.CreateEmojiData) (*emoji.Emoji, error) {
	return rest.CreateGuildEmoji(s.Token, guildId, data)
}

// updating Image is not permitted
func (s *Shard) ModifyGuildEmoji(guildId, emojiId uint64, data rest.CreateEmojiData) (*emoji.Emoji, error) {
	return rest.ModifyGuildEmoji(s.Token, guildId, emojiId, data)
}

func (s *Shard) CreateGuild(data rest.CreateGuildData) (*guild.Guild, error) {
	return rest.CreateGuild(s.Token, data)
}

func (s *Shard) GetGuild(guildId uint64) (*guild.Guild, error) {
	shouldCache := (*s.Cache).GetOptions().Guilds

	if shouldCache {
		cachedGuild := (*s.Cache).GetGuild(guildId)
		if cachedGuild != nil {
			return cachedGuild, nil
		}
	}

	guild, err := rest.GetGuild(s.Token, guildId)
	if err == nil {
		lock := (*s.Cache).GetLock(guildId)
		lock.Lock()
		(*s.Cache).StoreGuild(guild)
		lock.Unlock()
	}

	return guild, err
}

func (s *Shard) GetGuildPreview(guildId uint64) (*guild.GuildPreview, error) {
	return rest.GetGuildPreview(s.Token, guildId)
}

func (s *Shard) ModifyGuild(guildId uint64, data rest.ModifyGuildData) (*guild.Guild, error) {
	return rest.ModifyGuild(s.Token, guildId, data)
}

func (s *Shard) DeleteGuild(guildId uint64) error {
	return rest.DeleteGuild(s.Token, guildId)
}

func (s *Shard) GetGuildChannels(guildId uint64) ([]*channel.Channel, error) {
	shouldCache := (*s.Cache).GetOptions().Guilds && (*s.Cache).GetOptions().Channels

	if shouldCache {
		guild := (*s.Cache).GetGuild(guildId)
		if guild != nil {
			return guild.Channels, nil
		}
	}

	channels, err := rest.GetGuildChannels(s.Token, guildId)

	if shouldCache && err == nil {
		lock := (*s.Cache).GetLock(guildId)
		lock.Lock()

		guild := (*s.Cache).GetGuild(guildId)
		guild.Channels = channels
		(*s.Cache).StoreGuild(guild)

		lock.Unlock()
	}

	return channels, err
}

func (s *Shard) CreateGuildChannel(guildId uint64, data rest.CreateChannelData) (*channel.Channel, error) {
	return rest.CreateGuildChannel(s.Token, guildId, data)
}

func (s *Shard) ModifyGuildChannelPositions(guildId uint64, positions []rest.Position) error {
	return rest.ModifyGuildChannelPositions(s.Token, guildId, positions)
}

func (s *Shard) GetGuildMember(guildId, userId uint64) (*member.Member, error) {
	cacheGuilds := (*s.Cache).GetOptions().Guilds
	cacheUsers := (*s.Cache).GetOptions().Users

	if cacheGuilds && cacheUsers {
		guild := (*s.Cache).GetGuild(guildId)
		if guild != nil {
			for _, member := range guild.Members {
				if member != nil {
					if member.User != nil {
						if member.User.Id == userId {
							return member, nil
						}
					}
				}
			}
		}
	}

	member, err := rest.GetGuildMember(s.Token, guildId, userId)

	if cacheGuilds && err == nil {
		go func() {
			lock := (*s.Cache).GetLock(guildId)
			lock.Lock()

			guild := (*s.Cache).GetGuild(guildId)
			if guild != nil {
				guild.Members = append(guild.Members, member)
			}
			(*s.Cache).StoreGuild(guild)

			lock.Unlock()
		}()
	}

	if cacheUsers && err == nil {
		go func() {
			lock := (*s.Cache).GetLock(userId)
			lock.Lock()

			(*s.Cache).StoreUser(member.User)

			lock.Unlock()
		}()
	}

	return member, err
}

func (s *Shard) ListGuildMembers(guildId uint64, data rest.ListGuildMembersData) ([]*member.Member, error) {
	members, err := rest.ListGuildMembers(s.Token, guildId, data)

	cacheGuilds := (*s.Cache).GetOptions().Guilds
	cacheUsers := (*s.Cache).GetOptions().Users

	if cacheGuilds && err == nil {
		go func() {
			lock := (*s.Cache).GetLock(guildId)
			lock.Lock()

			guild := (*s.Cache).GetGuild(guildId)
			if guild != nil {
				new := make([]*member.Member, 0)

				for _, retrieved := range members {
					found := false

					internal:
					for _, existing := range guild.Members {
						if retrieved.User.Id == existing.User.Id {
							found = true
							break internal
						}
					}

					if !found {
						new = append(new, retrieved)
					}
				}

				guild.Members = append(guild.Members, new...)
			}
			(*s.Cache).StoreGuild(guild)

			lock.Unlock()
		}()
	}

	if cacheUsers && err == nil {
		go func() {
			for _, member := range members {
				lock := (*s.Cache).GetLock(member.User.Id)
				lock.Lock()

				(*s.Cache).StoreUser(member.User)

				lock.Unlock()
			}
		}()
	}

	return members, err
}

func (s *Shard) ModifyGuildMember(guildId, userId uint64, data rest.ModifyGuildMemberData) error {
	return rest.ModifyGuildMember(s.Token, guildId, userId, data)
}

func (s *Shard) ModifyCurrentUserNick(guildId uint64, nick string) error {
	return rest.ModifyCurrentUserNick(s.Token, guildId, nick)
}

func (s *Shard) AddGuildMemberRole(guildId, userId, roleId uint64) error {
	return rest.AddGuildMemberRole(s.Token, guildId, userId, roleId)
}

func (s *Shard) RemoveGuildMemberRole(guildId, userId, roleId uint64) error {
	return rest.RemoveGuildMemberRole(s.Token, guildId, userId, roleId)
}

func (s *Shard) RemoveGuildMember(guildId, userId uint64) error {
	return rest.RemoveGuildMember(s.Token, guildId, userId)
}

func (s *Shard) GetGuildBans(guildId uint64) ([]*guild.Ban, error) {
	return rest.GetGuildBans(s.Token, guildId)
}

func (s *Shard) GetGuildBan(guildId, userId uint64) (*guild.Ban, error) {
	return rest.GetGuildBan(s.Token, guildId, userId)
}

func (s *Shard) CreateGuildBan(guildId, userId uint64, data rest.CreateGuildBanData) error {
	return rest.CreateGuildBan(s.Token, guildId, userId, data)
}

func (s *Shard) RemoveGuildBan(guildId, userId uint64) error {
	return rest.RemoveGuildBan(s.Token, guildId, userId)
}

func (s *Shard) GetGuildRoles(guildId uint64) ([]*guild.Role, error) {
	shouldCache := (*s.Cache).GetOptions().Guilds

	if shouldCache {
		cachedGuild := (*s.Cache).GetGuild(guildId)
		if cachedGuild != nil {
			return cachedGuild.Roles, nil
		}
	}

	roles, err := rest.GetGuildRoles(s.Token, guildId)

	if shouldCache && err == nil {
		go func() {
			lock := (*s.Cache).GetLock(guildId)
			lock.Lock()

			cachedGuild := (*s.Cache).GetGuild(guildId)
			if cachedGuild == nil {
				cachedGuild = &guild.Guild{
					Id: guildId,
				}
			}
			cachedGuild.Roles = roles
			(*s.Cache).StoreGuild(cachedGuild)
			lock.Unlock()
		}()
	}

	return roles, err
}

func (s *Shard) CreateGuildRole(guildId uint64, data rest.GuildRoleData) (*guild.Role, error) {
	return rest.CreateGuildRole(s.Token, guildId, data)
}

func (s *Shard) ModifyGuildRolePositions(guildId uint64, positions []rest.Position) ([]*guild.Role, error) {
	return rest.ModifyGuildRolePositions(s.Token, guildId, positions)
}

func (s *Shard) ModifyGuildRole(guildId, roleId uint64, data rest.GuildRoleData) (*guild.Role, error) {
	return rest.ModifyGuildRole(s.Token, guildId, roleId, data)
}

func (s *Shard) DeleteGuildRole(guildId, roleId uint64) error {
	return rest.DeleteGuildRole(s.Token, guildId, roleId)
}

func (s *Shard) GetGuildPruneCount(guildId uint64, days int) (int, error) {
	return rest.GetGuildPruneCount(s.Token, guildId, days)
}

// computePruneCount = whether 'pruned' is returned, discouraged for large guilds
func (s *Shard) BeginGuildPrune(guildId uint64, days int, computePruneCount bool) error {
	return rest.BeginGuildPrune(s.Token, guildId, days, computePruneCount)
}

func (s *Shard) GetGuildVoiceRegions(guildId uint64) ([]*guild.VoiceRegion, error) {
	return rest.GetGuildVoiceRegions(s.Token, guildId)
}

func (s *Shard) GetGuildInvites(guildId uint64) ([]*invite.InviteMetadata, error) {
	return rest.GetGuildInvites(s.Token, guildId)
}

func (s *Shard) GetGuildIntegrations(guildId uint64) ([]*integration.Integration, error) {
	return rest.GetGuildIntegrations(s.Token, guildId)
}

func (s *Shard) CreateGuildIntegration(guildId uint64, data rest.CreateIntegrationData) error {
	return rest.CreateGuildIntegration(s.Token, guildId, data)
}

func (s *Shard) ModifyGuildIntegration(guildId, integrationId uint64, data rest.ModifyIntegrationData) error {
	return rest.ModifyGuildIntegration(s.Token, guildId, integrationId, data)
}

func (s *Shard) DeleteGuildIntegration(guildId, integrationId uint64) error {
	return rest.DeleteGuildIntegration(s.Token, guildId, integrationId)
}

func (s *Shard) SyncGuildIntegration(guildId, integrationId uint64) error {
	return rest.SyncGuildIntegration(s.Token, guildId, integrationId)
}

func (s *Shard) GetGuildEmbed(guildId uint64) (*guild.GuildEmbed, error) {
	return rest.GetGuildEmbed(s.Token, guildId)
}

func (s *Shard) ModifyGuildEmbed(guildId uint64, data guild.GuildEmbed) (*guild.GuildEmbed, error) {
	return rest.ModifyGuildEmbed(s.Token, guildId, data)
}

// returns invite object with only "code" and "uses" fields
func (s *Shard) GetGuildVanityUrl(guildId uint64) (*invite.Invite, error) {
	return rest.GetGuildVanityURL(s.Token, guildId)
}

func (s *Shard) GetGuildWidgetImage(guildId uint64, style guild.WidgetStyle) (*image.Image, error) {
	return rest.GetGuildWidgetImage(s.Token, guildId, style)
}

func (s *Shard) GetInvite(inviteCode string, withCounts bool) (*invite.Invite, error) {
	return rest.GetInvite(s.Token, inviteCode, withCounts)
}

func (s *Shard) DeleteInvite(inviteCode string) (*invite.Invite, error) {
	return rest.DeleteInvite(s.Token, inviteCode)
}

func (s *Shard) GetCurrentUser() (*user.User, error) {
	if cached := (*s.Cache).GetSelf(); cached != nil {
		return cached, nil
	}

	self, err := rest.GetCurrentUser(s.Token)

	if err == nil {
		go func() {
			(*s.Cache).StoreSelf(self)
		}()
	}

	return self, err
}

func (s *Shard) GetUser(userId uint64) (*user.User, error) {
	shouldCache := (*s.Cache).GetOptions().Users

	if shouldCache {
		cached := (*s.Cache).GetUser(userId)
		if cached != nil {
			return cached, nil
		}
	}

	user, err := rest.GetUser(s.Token, userId)

	if shouldCache && err == nil {
		go func() {
			lock := (*s.Cache).GetLock(userId)
			lock.Lock()

			(*s.Cache).StoreUser(user)

			lock.Unlock()
		}()
	}

	return user, err
}

func (s *Shard) ModifyCurrentUser(data rest.ModifyUserData) (*user.User, error) {
	return rest.ModifyCurrentUser(s.Token, data)
}

func (s *Shard) GetCurrentUserGuilds(data rest.CurrentUserGuildsData) ([]*guild.Guild, error) {
	return rest.GetCurrentUserGuilds(s.Token, data)
}

func (s *Shard) LeaveGuild(guildId uint64) error {
	return rest.LeaveGuild(s.Token, guildId)
}

func (s *Shard) CreateDM(recipientId uint64) (*channel.Channel, error) {
	return rest.CreateDM(s.Token, recipientId)
}

func (s *Shard) GetUserConnections() ([]*integration.Connection, error) {
	return rest.GetUserConnections(s.Token)
}

// GetGuildVoiceRegions should be preferred, as it returns VIP servers if available to the guild
func (s *Shard) ListVoiceRegions() ([]*guild.VoiceRegion, error) {
	return rest.ListVoiceRegions(s.Token)
}

func (s *Shard) CreateWebhook(channelId uint64, data rest.WebhookData) (*guild.Webhook, error) {
	return rest.CreateWebhook(s.Token, channelId, data)
}

func (s *Shard) GetChannelWebhooks(channelId uint64) ([]*guild.Webhook, error) {
	return rest.GetChannelWebhooks(s.Token, channelId)
}

func (s *Shard) GetGuildWebhooks(guildId uint64) ([]*guild.Webhook, error) {
	return rest.GetGuildWebhooks(s.Token, guildId)
}

func (s *Shard) GetWebhook(webhookId uint64) (*guild.Webhook, error) {
	return rest.GetWebhook(s.Token, webhookId)
}

func (s *Shard) ModifyWebhook(webhookId uint64, data rest.ModifyWebhookData) (*guild.Webhook, error) {
	return rest.ModifyWebhook(s.Token, webhookId, data)
}

func (s *Shard) DeleteWebhook(webhookId uint64) error {
	return rest.DeleteWebhook(s.Token, webhookId)
}

func (s *Shard) ExecuteWebhook(webhookId uint64, webhookToken string, wait bool, data rest.WebhookBody) {
	rest.ExecuteWebhook(webhookId, webhookToken, wait, data)
}
