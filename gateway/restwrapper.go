package gateway

import (
	"github.com/Dot-Rar/gdl/objects"
	"github.com/Dot-Rar/gdl/rest"
	"github.com/sirupsen/logrus"
)

func (s *Shard) GetChannel(channelId uint64) *objects.Channel {
	cacheChannels := (*s.Cache).GetOptions().Channels
	if cacheChannels {
		cached := (*s.Cache).GetChannel(channelId)
		if cached != nil {
			return cached
		}
	}

	channel, err := rest.GetChannel(channelId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing GetChannel: %s", err.Error())
		return nil
	}

	if cacheChannels {
		(*s.Cache).StoreChannel(channel)
	}

	return channel
}

func (s *Shard) ModifyChannel(channelId uint64, data rest.ModifyChannelData) *objects.Channel {
	channel, err := rest.ModifyChannel(channelId, s.Token, data)
	if err != nil {
		logrus.Warnf("error while executing ModifyChannel: %s", err.Error())
		return nil
	}

	if (*s.Cache).GetOptions().Channels {
		(*s.Cache).StoreChannel(channel)
	}

	return channel
}

func (s *Shard) DeleteChannel(channelId uint64) *objects.Channel {
	channel, err := rest.DeleteChannel(channelId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing DeleteChannel: %s", err.Error())
		return nil
	}

	if (*s.Cache).GetOptions().Channels {
		(*s.Cache).DeleteChannel(channelId)
	}

	return channel
}

func (s *Shard) GetChannelMessages(channelId uint64, options rest.GetChannelMessagesData) []objects.Message {
	messages, err := rest.GetChannelMessages(channelId, s.Token, options)
	if err != nil {
		logrus.Warnf("error while executing GetChannelMessages: %s", err.Error())
		return make([]objects.Message, 0)
	}

	return messages
}

func (s *Shard) GetChannelMessage(channelId, messageId uint64) *objects.Message {
	message, err := rest.GetChannelMessage(channelId, messageId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing GetChannelMessage: %s", err.Error())
		return nil
	}

	return message
}

func (s *Shard) CreateMessage(channelId uint64, content string) *objects.Message {
	return s.CreateMessageComplex(channelId, rest.CreateMessageData{
		Content: content,
	})
}

func (s *Shard) CreateMessageComplex(channelId uint64, data rest.CreateMessageData) *objects.Message {
	message, err := rest.CreateMessage(channelId, s.Token, data)
	if err != nil {
		logrus.Warnf("error while executing CreateMessage: %s", err.Error())
		return nil
	}

	return message
}

func (s *Shard) CreateReaction(channelId, messageId uint64, emoji string) {
	err := rest.CreateReaction(channelId, messageId, emoji, s.Token)
	if err != nil {
		logrus.Warnf("error while executing CreateReaction: %s", err.Error())
	}
}

func (s *Shard) DeleteOwnReaction(channelId, messageId uint64, emoji string) {
	err := rest.DeleteOwnReaction(channelId, messageId, emoji, s.Token)
	if err != nil {
		logrus.Warnf("error while executing DeleteOwnReaction: %s", err.Error())
	}
}

func (s *Shard) DeleteUserReaction(channelId, messageId, userId uint64, emoji string) {
	err := rest.DeleteUserReaction(channelId, messageId, userId, emoji, s.Token)
	if err != nil {
		logrus.Warnf("error while executing DeleteUserReaction: %s", err.Error())
	}
}

func (s *Shard) GetReactions(channelId, messageId uint64, emoji string, options rest.GetReactionsData) []objects.User {
	users, err := rest.GetReactions(channelId, messageId, emoji, s.Token, options)
	if err != nil {
		logrus.Warnf("error while executing GetReactions: %s", err.Error())
		return make([]objects.User, 0)
	}

	return users
}

func (s *Shard) DeleteAllReactions(channelId, messageId uint64) {
	err := rest.DeleteAllReactions(channelId, messageId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing DeleteAllReactions: %s", err.Error())
	}
}

func (s *Shard) DeleteAllReactionsEmoji(channelId, messageId uint64, emoji string) {
	err := rest.DeleteAllReactionsEmoji(channelId, messageId, emoji, s.Token)
	if err != nil {
		logrus.Warnf("error while executing DeleteAllReactionsEmoji: %s", err.Error())
	}
}

func (s *Shard) EditMessage(channelId, messageId uint64, data rest.ModifyChannelData) *objects.Message {
	message, err := rest.EditMessage(channelId, messageId, s.Token, data)
	if err != nil {
		logrus.Warnf("error while executing EditMessage: %s", err.Error())
		return nil
	}

	return message
}

func (s *Shard) DeleteMessage(channelId, messageId uint64) {
	err := rest.DeleteMessage(channelId, messageId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing DeleteMessage: %s", err.Error())
	}
}

func (s *Shard) BulkDeleteMessages(channelId uint64, messages []uint64) {
	err := rest.BulkDeleteMessages(channelId, messages, s.Token)
	if err != nil {
		logrus.Warnf("error while executing BulkDeleteMessages: %s", err.Error())
	}
}

func (s *Shard) EditChannelPermissions(channelId uint64, updated objects.Overwrite) {
	err := rest.EditChannelPermissions(channelId, s.Token, updated)
	if err != nil {
		logrus.Warnf("error while executing EditChannelPermissions: %s", err.Error())
	}
}

func (s *Shard) GetChannelInvites(channelId uint64) []objects.InviteMetadata {
	invites, err := rest.GetChannelInvites(channelId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing GetChannelInvites: %s", err.Error())
		return make([]objects.InviteMetadata, 0)
	}

	return invites
}

func (s *Shard) CreateChannelInvite(channelId uint64, data objects.InviteMetadata) *objects.Invite {
	invite, err := rest.CreateChannelInvite(channelId, s.Token, data)
	if err != nil {
		logrus.Warnf("error while executing CreateChannelInvite: %s", err.Error())
		return nil
	}

	return invite
}

func (s *Shard) DeleteChannelPermissions(channelId, overwriteId uint64) {
	err := rest.DeleteChannelPermissions(channelId, overwriteId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing DeleteChannelPermissions: %s", err.Error())
	}
}

func (s *Shard) TriggerTypingIndicator(channelId uint64) {
	err := rest.TriggerTypingIndicator(channelId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing TriggerTypingIndicator: %s", err.Error())
	}
}

func (s *Shard) GetPinnedMessages(channelId uint64) []objects.Message {
	messages, err := rest.GetPinnedMessages(channelId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing GetPinnedMessages: %s", err.Error())
		return make([]objects.Message, 0)
	}

	return messages
}

func (s *Shard) AddPinnedChannelMessage(channelId, messageId uint64) {
	err := rest.AddPinnedChannelMessage(channelId, messageId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing AddPinnedChannelMessage: %s", err.Error())
	}
}

func (s *Shard) DeletePinnedChannelMessage(channelId, messageId uint64) {
	err := rest.DeletePinnedChannelMessage(channelId, messageId, s.Token)
	if err != nil {
		logrus.Warnf("error while executing DeletePinnedChannelMessage: %s", err.Error())
	}
}
