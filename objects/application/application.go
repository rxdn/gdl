package application

import (
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/user"
)

type Application struct {
	Id                             uint64         `json:"id,string"`
	Name                           string         `json:"name"`
	Icon                           *string        `json:"icon"`
	Description                    string         `json:"description"`
	RpcOrigins                     []string       `json:"rpc_origins,omitempty"`
	BotPublic                      bool           `json:"bot_public"`
	BotRequireCodeGrant            bool           `json:"bot_require_code_grant"`
	Bot                            *user.User     `json:"bot,omitempty"`
	TermsOfServiceUrl              *string        `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyUrl               *string        `json:"privacy_policy_url,omitempty"`
	Owner                          *user.User     `json:"owner,omitempty"`
	VerifyKey                      string         `json:"verify_key"`
	Team                           *Team          `json:"team"`
	GuildId                        *uint64        `json:"guild_id,string,omitempty"`
	Guild                          *guild.Guild   `json:"guild,omitempty"`
	PrimarySkuId                   *uint64        `json:"primary_sku_id,string,omitempty"`
	Slug                           *string        `json:"slug,omitempty"`
	CoverImage                     *string        `json:"cover_image,omitempty"`
	Flags                          *Flag          `json:"flags,omitempty"`
	ApproximateGuildCount          *int           `json:"approximate_guild_count,omitempty"`
	RedirectUris                   []string       `json:"redirect_uris,omitempty"`
	InteractionsEndpointUrl        *string        `json:"interactions_endpoint_url,omitempty"`
	RoleConnectionsVerificationUrl *string        `json:"role_connections_verification_url,omitempty"`
	Tags                           []string       `json:"tags,omitempty"`
	InstallParams                  *InstallParams `json:"install_params,omitempty"`
	CustomInstallUrl               *string        `json:"custom_install_url,omitempty"`
}

type Flag uint64

const (
	FlagAutoModerationBadge              Flag = 1 << 6
	FlagIntentGatewayPresence            Flag = 1 << 12
	FlagIntentGatewayPresenceLimited     Flag = 1 << 13
	FlagIntentGatewayGuildMembers        Flag = 1 << 14
	FlagIntentGatewayGuildMembersLimited Flag = 1 << 15
	FlagVerificationPendingGuildLimit    Flag = 1 << 16
	FlagEmbedded                         Flag = 1 << 17
	FlagGatewayMessageContent            Flag = 1 << 18
	FlagGatewayMessageContentLimited     Flag = 1 << 19
	FlagApplicationCommandBadge          Flag = 1 << 20
)

func (f Flag) Has(flag Flag) bool {
	return f&flag == flag
}

func BuildFlags(flags ...Flag) Flag {
	var built Flag

	for _, flag := range flags {
		built |= flag
	}

	return built
}
