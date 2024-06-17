package rest

import (
	"context"
	"github.com/rxdn/gdl/objects/application"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
)

func GetCurrentApplication(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter) (application.Application, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    "/applications/@me",
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteGetCurrentApplication, 0),
		RateLimiter: rateLimiter,
	}

	var app application.Application
	err, _ := endpoint.Request(ctx, token, nil, &app)
	return app, err
}

type EditCurrentApplicationData struct {
	CustomInstallUrl               *string                    `json:"custom_install_url,omitempty"`
	Description                    *string                    `json:"description,omitempty"`
	RoleConnectionsVerificationUrl *string                    `json:"role_connections_verification_url,omitempty"`
	InstallParams                  *application.InstallParams `json:"install_params,omitempty"`
	Flags                          *application.Flag          `json:"flags,omitempty"`
	// TODO: icon
	// TODO: cover_image
	InteractionsEndpointUrl *string  `json:"interactions_endpoint_url,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
}

func EditCurrentApplication(ctx context.Context, token string, rateLimiter *ratelimit.Ratelimiter, data EditCurrentApplicationData) (application.Application, error) {
	endpoint := request.Endpoint{
		RequestType: request.PATCH,
		ContentType: request.ApplicationJson,
		Endpoint:    "/applications/@me",
		Route:       ratelimit.NewApplicationRoute(ratelimit.RouteEditCurrentApplication, 0),
		RateLimiter: rateLimiter,
	}

	var app application.Application
	err, _ := endpoint.Request(ctx, token, data, &app)
	return app, err
}
