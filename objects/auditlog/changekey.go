package auditlog

type ChangeKey string

const (
	// guild
	ChangeKeyName                        ChangeKey = "name"
	ChangeKeyIconHash                    ChangeKey = "icon_hash"
	ChangeKeySplashHash                  ChangeKey = "splash_hash"
	ChangeKeyOwnerId                     ChangeKey = "owner_id"
	ChangeKeyRegion                      ChangeKey = "region"
	ChangeKeyAfkChannelId                ChangeKey = "afk_channel_id"
	ChangeKeyAfkTimeout                  ChangeKey = "afk_timeout"
	ChangeKeyMfaLevel                    ChangeKey = "mfa_level"
	ChangeKeyVerificationLevel           ChangeKey = "verification_level"
	ChangeKeyExplicitContentFilter       ChangeKey = "explicit_content_filter"
	ChangeKeyDefaultMessageNotifications ChangeKey = "default_message_notifications"
	ChangeKeyVanityUrlCode               ChangeKey = "vanity_url_code"
	ChangeKeyRoleAdd                     ChangeKey = "$add"
	ChangeKeyRoleRemove                  ChangeKey = "$remove"
	ChangeKeyPruneDeleteDays             ChangeKey = "prune_delete_days"
	ChangeKeyWidgetEnabled               ChangeKey = "widget_enabled"
	ChangeKeyWidgetChannelId             ChangeKey = "widget_channel_id"
	ChangeKeySystemChannelId             ChangeKey = "system_channel_id"

	// channel
	ChangeKeyPosition             ChangeKey = "position"
	ChangeKeyTopic                ChangeKey = "topic"
	ChangeKeyBitrate              ChangeKey = "bitrate"
	ChangeKeyPermissionOverwrites ChangeKey = "permission_overwrites"
	ChangeKeyNsfw                 ChangeKey = "nsfw"
	ChangeKeyApplicationId        ChangeKey = "application_id"
	ChangeKeyRateLimitPerUser     ChangeKey = "rate_limit_per_user"

	// role
	ChangeKeyPermissions ChangeKey = "permissions"
	ChangeKeyColor       ChangeKey = "color"
	ChangeKeyHoist       ChangeKey = "hoist"
	ChangeKeyMentionable ChangeKey = "mentionable"
	ChangeKeyAllow       ChangeKey = "allow"
	ChangeKeyDeny        ChangeKey = "deny"

	// invite
	ChangeKeyCode      ChangeKey = "code"
	ChangeKeyChannelId ChangeKey = "channel_id"
	ChangeKeyInviterId ChangeKey = "inviter_id"
	ChangeKeyMaxUses   ChangeKey = "max_uses"
	ChangeKeyUses      ChangeKey = "uses"
	ChangeKeyMaxAge    ChangeKey = "max_age"
	ChangeKeyTemporary ChangeKey = "temporary"

	// user
	ChangeKeyDeaf       ChangeKey = "deaf"
	ChangeKeyMute       ChangeKey = "mute"
	ChangeKeyNick       ChangeKey = "nick"
	ChangeKeyAvatarHash ChangeKey = "avatar_hash"

	// any
	ChangeKeyId   ChangeKey = "id"
	ChangeKeyType ChangeKey = "type"

	// integration
	ChangeKeyEnableEmoticons   ChangeKey = "enable_emoticons"
	ChangeKeyExpireBehaviour   ChangeKey = "expire_behaviour"
	ChangeKeyExpireGracePeriod ChangeKey = "expire_grace_period"
)
