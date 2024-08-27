package entitlement

import "time"

type Entitlement struct {
	Id            uint64          `json:"id,string"`
	SkuId         uint64          `json:"sku_id,string"`
	ApplicationId uint64          `json:"application_id,string"`
	UserId        *uint64         `json:"user_id,string,omitempty"`
	Type          EntitlementType `json:"type"`
	Deleted       bool            `json:"deleted"`
	StartsAt      *time.Time      `json:"starts_at,omitempty"`
	EndsAt        *time.Time      `json:"ends_at,omitempty"`
	GuildId       *uint64         `json:"guild_id,string,omitempty"`
	Consumed      *bool           `json:"consumed,omitempty"`
}

type EntitlementType uint16

const (
	TypePurchase EntitlementType = iota + 1
	TypePremiumSubscription
	TypeDeveloperGift
	TypeTestModePurchase
	TypeFreePurchase
	TypeUserGift
	TypePremiumPurchase
	TypeApplicationSubscription
)
