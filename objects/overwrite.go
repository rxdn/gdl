package objects

import (
	"fmt"
	"github.com/Dot-Rar/gdl/utils"
)

type Overwrite struct {
	Id    uint64 `json:"id,string"`
	Type  string `json:"type"`
	Allow int    `json:"allow"`
	Deny  int    `json:"deny"`
}

func (o *Overwrite) KeyName() string {
	return fmt.Sprintf("cache:overwrite:%s", o.Id)
}

func (o *Overwrite) Serialize() map[string]map[string]interface{} {
	fields := make(map[string]map[string]interface{})
	utils.CopyNonNil(fields, o.KeyName(), o)
	return fields
}
