package rest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
)

type (
	GrantType string

	TokenExchangeBody struct {
		GrantType    GrantType `qs:"grant_type"`
		RefreshToken string    `qs:"refresh_token,omitempty"`
		Code         string    `qs:"code,omitempty"`
		RedirectUri  string    `qs:"redirect_uri,omitempty"`
	}

	TokenExchangeResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	TokenTypeHint string

	TokenRevocationBody struct {
		Token         string         `qs:"token"`
		TokenTypeHint *TokenTypeHint `qs:"token_type_hint,omitempty"`
	}
)

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"

	TokenTypeHintAccessToken  TokenTypeHint = "access_token"
	TokenTypeHintRefreshToken TokenTypeHint = "refresh_token"
)

func ExchangeCode(ctx context.Context, rateLimiter *ratelimit.Ratelimiter, clientId uint64, clientSecret, redirectUri, code string) (TokenExchangeResponse, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationFormUrlEncoded,
		Endpoint:    "/oauth2/token",
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteOauth2TokenExchange, clientId),
		RateLimiter: rateLimiter,
	}

	body := TokenExchangeBody{
		GrantType:   GrantTypeAuthorizationCode,
		Code:        code,
		RedirectUri: redirectUri,
	}

	header := "Basic " + base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%s", clientId, clientSecret)))

	var tokenExchangeResponse TokenExchangeResponse
	err, _ := endpoint.Request(ctx, header, body, &tokenExchangeResponse)
	return tokenExchangeResponse, convertToOauthError(err)
}

func RefreshToken(ctx context.Context, rateLimiter *ratelimit.Ratelimiter, clientId uint64, clientSecret, refreshToken string) (TokenExchangeResponse, error) {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationFormUrlEncoded,
		Endpoint:    "/oauth2/token",
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteOauth2TokenExchange, clientId),
		RateLimiter: rateLimiter,
	}

	body := TokenExchangeBody{
		GrantType:    GrantTypeRefreshToken,
		RefreshToken: refreshToken,
	}

	var tokenExchangeResponse TokenExchangeResponse
	header := "Basic " + base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%s", clientId, clientSecret)))
	err, _ := endpoint.Request(ctx, header, body, &tokenExchangeResponse)
	return tokenExchangeResponse, convertToOauthError(err)
}

func RevokeToken(ctx context.Context, rateLimiter *ratelimit.Ratelimiter, clientId uint64, clientSecret, token string, tokenTypeHint *TokenTypeHint) error {
	endpoint := request.Endpoint{
		RequestType: request.POST,
		ContentType: request.ApplicationFormUrlEncoded,
		Endpoint:    "/oauth2/token/revoke",
		Route:       ratelimit.NewOtherRoute(ratelimit.RouteOauth2TokenRevoke, clientId),
		RateLimiter: rateLimiter,
	}

	body := TokenRevocationBody{
		Token:         token,
		TokenTypeHint: tokenTypeHint,
	}

	header := "Basic " + base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%s", clientId, clientSecret)))
	err, _ := endpoint.Request(ctx, header, body, nil)
	return convertToOauthError(err)
}

func convertToOauthError(err error) error {
	var restError request.RestError
	if errors.As(err, &restError) {
		var oauthError request.OAuthError
		if err := json.Unmarshal(restError.Raw, &oauthError); err != nil {
			return fmt.Errorf("error deserialize oauth error: %w", err)
		}

		return oauthError
	}

	return err
}
