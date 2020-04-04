package guild

type GuildFeature string

const (
	GuildFeatureInviteSplash   GuildFeature = "INVITE_SPLASH"
	GuildFeatureVipRegions     GuildFeature = "VIP_REGIONS" // guild has access to set 384kbps bitrate in voice (previously VIP voice servers)
	GuildFeatureVanityUrl      GuildFeature = "VANITY_URL"
	GuildFeatureVerified       GuildFeature = "VERIFIED"
	GuildFeaturePartnered      GuildFeature = "PARTNERED"
	GuildFeaturePublic         GuildFeature = "PUBLIC"
	GuildFeatureCommerce       GuildFeature = "COMMERCE"
	GuildFeatureNews           GuildFeature = "NEWS"
	GuildFeatureDiscoverable   GuildFeature = "DISCOVERABLE"
	GuildFeatureFeaturable     GuildFeature = "FEATURABLE"
	GuildFeatureBanner         GuildFeature = "BANNER"
	GuildFeaturePublicDisabled GuildFeature = "PUBLIC_DISABLED"
)
