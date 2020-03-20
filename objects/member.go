package objects

import (
	"fmt"
	"github.com/Dot-Rar/gdl/utils"
	"github.com/fatih/structs"
	"reflect"
)

type Member struct {
	User     *User     `json:"user"`
	Nick     string   `json:"nick"`
	Roles    []uint64 `json:"roles,string"`
	JoinedAt string   `json:"joined_at"`
	Deaf     bool     `json:"deaf"`
	Mute     bool     `json:"mute"`
}

func (m *Member) KeyName(guildId string) string {
	return fmt.Sprintf("cache:member:%s:%s", guildId, m.User.Id)
}

func (m *Member) Serialize(guildId string) map[string]map[string]interface{} {
	fields := make(map[string]map[string]interface{})
	utils.Initialise(fields, m.KeyName(guildId))

	for k, v := range structs.Map(m) {
		if k == "User" {
			fields = utils.Append(fields, m.User.Serialize())
			fields[m.KeyName(guildId)][k] = m.User.Id
		} else { // This handles the Roles field fine too
			if !utils.IsZero(reflect.ValueOf(v)) {
				fields[m.KeyName(guildId)][k] = v
			}
		}
	}

	return fields
}
