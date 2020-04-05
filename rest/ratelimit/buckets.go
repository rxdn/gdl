package ratelimit

import (
	"fmt"
	"strconv"
)

// https://github.com/discord/discord-api-docs/issues/981#issuecomment-534634946
// cake is really smart

func NewGatewayBucket() string {
	return "41f9cd5d28af77da04563bcb1d67fdfd"
}

func NewMessageWriteBucket(channelId uint64) string {
	return fmt.Sprintf("80c17d2f203122d936070c88c8d10f33-%d", channelId)
}

func NewMemberModifyBucket(guildId uint64) string {
	return fmt.Sprintf("e670de274720e74c12bd7b9cfc111470-%d", guildId)
}

func NewGuildMemberModifyBucket(guildId uint64) string {
	return fmt.Sprintf("e1e879dfc5003bd7ca461c93344aed2f-%d", guildId)
}

func NewInviteCreateBucket(channelId uint64) string {
	return fmt.Sprintf("ad8bc58100d4c04b6352b432fc194883-%d", channelId)
}

func NewTypingBucket(channelId uint64) string {
	return fmt.Sprintf("09e618e6deafd34f65e8a56d13bb6c1e-%d", channelId)
}

func NewRoleCreateBucket(guildId uint64) string {
	return fmt.Sprintf("019dc22942a45e5a5889485b28769a18-%d", guildId)
}

func NewLobbiesBucket() string {
	return "4bf541531275fa822fb3b6ea6cbedd75"
}

func NewReactionsBucket(channelId uint64) string {
	return fmt.Sprintf("d6b7697c78814ecd72bb20df05517c78-%d", channelId)
}

func NewLobbiesWriteBucket(lobbyId uint64) string {
	return fmt.Sprintf("d8671d1cb097a3bd26b2344ee25a0562-%d", lobbyId)
}

func NewModifyCurrentNickBucket(guildId uint64) string {
	return fmt.Sprintf("11fab958fe94bf2a06661d856771b522-%d", guildId)
}

func NewWebhookExecuteBucket(webhookId uint64) string {
	return fmt.Sprintf("3cd1f278bd0ecaf11e0d2391374c011d-%d", webhookId)
}

func NewGuildListMembersBucket(guildId uint64) string {
	return fmt.Sprintf("db91b7148e84e53dbe2bdb5f6158bddf-%d", guildId)
}

// i need to get more creative with these names
func NewGuildMemberGetOrDeleteBucket(guildId uint64) string {
	return fmt.Sprintf("63e8b9f22df9dc9ef04cd65af6244664-%d", guildId)
}

func NewInviteGetBucket(code string) string {
	return fmt.Sprintf("5dddb913eea0498009a1d7f5a06a1643-%s", code)
}

func NewMessageDeleteBucket(channelId uint64) string {
	return fmt.Sprintf("087226e88721bc988cf853c666255256-%d", channelId)
}

func NewBulkDeleteBucket(channelId uint64) string {
	return fmt.Sprintf("b05c0d8c2ab83895085006a8eae073a3-%d", channelId)
}

// Generic buckets
func NewChannelBucket(channelId uint64) string {
	return strconv.FormatUint(channelId, 10)
}

func NewGuildBucket(guildId uint64) string {
	return strconv.FormatUint(guildId, 10)
}

// TODO: Investigate best way to handle this
func NewEmojiBucket(guildId uint64) string {
	return fmt.Sprintf("e-%d", guildId)
}

func NewUserBucket(id uint64) string {
	return strconv.FormatUint(id, 10)
}

func NewWebhookBucket(webhookId uint64) string {
	return fmt.Sprintf("/webhooks/%d", webhookId)
}
