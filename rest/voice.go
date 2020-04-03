package rest

import (
	"fmt"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/request"
	"github.com/rxdn/gdl/rest/routes"
)

func ListVoiceRegions(token string) ([]guild.VoiceRegion, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/voice/regions"),
	}

	var voiceRegions []guild.VoiceRegion
	err, _ := endpoint.Request(token, &routes.RouteManager.GetVoiceRoute().Ratelimiter, nil, &voiceRegions)
	return voiceRegions, err
}
