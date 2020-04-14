package payloads

import (
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/objects/user"
)

type (
	Identify struct {
		Opcode int          `json:"op"`
		Data   IdentifyData `json:"d"`
	}

	IdentifyData struct {
		Token              string            `json:"token"`
		Properties         Properties        `json:"properties"`
		Compress           bool              `json:"compress"`
		LargeThreshold     int               `json:"large_threshold"`
		Shard              []int             `json:"shard"`
		Presence           user.UpdateStatus `json:"presence"`
		GuildSubscriptions bool              `json:"guild_subscriptions"`
		Intents            int               `json:"intents,omitempty"`
	}

	Properties struct {
		Os      string `json:"$os"`
		Browser string `json:"$browser"`
		Device  string `json:"$device"`
	}

	Presence struct {
		Game   Game   `json:"game"`
		Status string `json:"status"`
		Since  *int   `json:"since"`
		Afk    bool   `json:"afk"`
	}

	Game struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	}
)

func NewIdentify(shardId int, shardTotal int, token string, status user.UpdateStatus, guildSubscriptions bool, intents []gateway.Intent) Identify {
	var sum int
	for _, intent := range intents {
		sum += int(intent)
	}

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
			Intents:            sum,
		},
	}

	return payload
}
