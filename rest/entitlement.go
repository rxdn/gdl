package rest

import (
	"context"
	"fmt"
	"github.com/rxdn/gdl/objects/entitlement"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"net/url"
	strconv "strconv"
)

type EntitlementQueryOptions struct {
	UserId        *uint64
	SkuIds        []uint64
	Before        *uint64
	After         *uint64
	Limit         *int
	GuildId       *uint64
	ExcludedEnded *bool
}

func (o *EntitlementQueryOptions) Query() string {
	query := url.Values{}

	if o.UserId != nil {
		query.Add("user_id", strconv.FormatUint(*o.UserId, 10))
	}

	if len(o.SkuIds) > 0 {
		var encoded string
		for i, id := range o.SkuIds {
			encoded += strconv.FormatUint(id, 10) + ","

			if i == len(o.SkuIds)-1 {
				encoded = encoded[:len(encoded)-1]
			}
		}

		query.Add("sku_ids", encoded)
	}

	if o.Before != nil {
		query.Add("before", strconv.FormatUint(*o.Before, 10))
	}

	if o.After != nil {
		query.Add("after", strconv.FormatUint(*o.After, 10))
	}

	if o.Limit != nil {
		query.Add("limit", strconv.Itoa(*o.Limit))
	}

	if o.GuildId != nil {
		query.Add("guild_id", strconv.FormatUint(*o.GuildId, 10))
	}

	if o.ExcludedEnded != nil {
		query.Add("excluded_ended", strconv.FormatBool(*o.ExcludedEnded))
	}

	return query.Encode()
}

func ListEntitlements(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64, options EntitlementQueryOptions) ([]entitlement.Entitlement, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/entitlements?%s", applicationId, options.Query()),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteListEntitlements, applicationId),
		RateLimiter: rateLimiter,
	}

	var entitlements []entitlement.Entitlement
	if err, _ := endpoint.Request(ctx, token, nil, &entitlements); err != nil {
		return nil, err
	}

	return entitlements, nil
}

func ConsumeEntitlement(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId, entitlementId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/entitlements/%d/consume", applicationId, entitlementId),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteConsumeEntitlement, applicationId),
		RateLimiter: rateLimiter,
	}

	if err, _ := endpoint.Request(ctx, token, nil, nil); err != nil {
		return err
	}

	return nil
}

type CreateTestEntitlementData struct {
	SkuId     uint64               `json:"sku_id,string"`
	OwnerId   uint64               `json:"owner_id,string"`
	OwnerType EntitlementOwnerType `json:"owner_type"`
}

type EntitlementOwnerType uint8

const (
	EntitlementOwnerTypeGuild EntitlementOwnerType = iota + 1
	EntitlementOwnerTypeUser
)

func CreateTestEntitlement(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, applicationId uint64, data CreateTestEntitlementData) (entitlement.Entitlement, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationJson,
		Endpoint:    fmt.Sprintf("/applications/%d/entitlements", applicationId),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteCreateTestEntitlement, applicationId),
		RateLimiter: rateLimiter,
	}

	var createdEntitlement entitlement.Entitlement
	if err, _ := endpoint.Request(ctx, token, data, &createdEntitlement); err != nil {
		return entitlement.Entitlement{}, err
	}

	return createdEntitlement, nil
}

func DeleteTestEntitlement(ctx context.Context, token string, rateLimtier *ratelimit.Ratelimiter, applicationId, entitlementId uint64) error {
	endpoint := request.Endpoint{
		RequestType: request.DELETE,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/applications/%d/entitlements/%d", applicationId, entitlementId),
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteDeleteTestEntitlement, applicationId),
		RateLimiter: rateLimtier,
	}

	err, _ := endpoint.Request(ctx, token, nil, nil)
	return err
}
