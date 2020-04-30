package payloads

import (
	gatewayintents "github.com/rxdn/gdl/gateway/intents"
	"github.com/rxdn/gdl/objects/user"
)

type (
	Identify struct {
		Opcode int          `json:"op"`
		Data   IdentifyData `json:"d"`
	}

	IdentifyData struct {
		Token              string     `json:"token"`
		Properties         Properties `json:"properties"`
		Compress           bool       `json:"compress"`
		LargeThreshold     int        `json:"large_threshold"`
		Shard              []int      `json:"shard"`
		Presence           Presence   `json:"presence,omitempty"`
		GuildSubscriptions bool       `json:"guild_subscriptions"`
		Intents            int        `json:"intents,omitempty"`
	}

	Properties struct {
		Os      string `json:"$os"`
		Browser string `json:"$browser"`
		Device  string `json:"$device"`
	}

	Presence struct {
		Game   Game   `json:"game,omitempty"`
		Status string `json:"status"`
		Since  *int   `json:"since,omitempty"`
		Afk    bool   `json:"afk"`
	}

	Game struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	}
)

func NewIdentify(shardId int, shardTotal int, token string, status user.UpdateStatus, guildSubscriptions bool, intents ...gatewayintents.Intent) Identify {
	payload := Identify{
		Opcode: 2,
		Data: IdentifyData{
			Token: token,
			Properties: Properties{
				Os:      "linux",
				Browser: "GDL",
				Device:  "GDL",
			},
			LargeThreshold:     250,
			Shard:              []int{shardId, shardTotal},
			Presence:           status,
			GuildSubscriptions: guildSubscriptions,
			Intents:            gatewayintents.SumIntents(intents...),
		},
	}

	return payload
}
