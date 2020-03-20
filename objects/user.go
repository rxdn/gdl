package objects

import (
	"fmt"
	"github.com/Dot-Rar/gdl/utils"
)

type User struct {
	Id            uint64 `json:"id,string"`
	Username      string `json:"username"`
	Discriminator uint16 `json:"discriminator,string"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	PremiumType   int    `json:"premium_type"`
}

func (u *User) KeyName() string {
	return fmt.Sprintf("cache:user:%s", u.Id)
}

func (u *User) Serialize() map[string]map[string]interface{} {
	fields := make(map[string]map[string]interface{})
	utils.CopyNonNil(fields, u.KeyName(), u)
	return fields
}
