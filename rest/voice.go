package rest

import (
	"fmt"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/request"
)

func ListVoiceRegions(token string) ([]guild.VoiceRegion, error) {
	endpoint := request.Endpoint{
		RequestType: request.GET,
		ContentType: request.Nil,
		Endpoint:    fmt.Sprintf("/voice/regions"),
	}

	var voiceRegions []guild.VoiceRegion
	err, _ := endpoint.Request(token, nil, &voiceRegions)
	return voiceRegions, err
}
