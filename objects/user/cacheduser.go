package user

import "fmt"

type CachedUser struct {
	Username    string  `json:"username"`
	GlobalName  *string `json:"global_name"`
	Avatar      string  `json:"avatar"`
	Bot         bool    `json:"bot"`
	Flags       uint64  `json:"flags"`
	PremiumType int     `json:"premium_type"`
}

func (u *CachedUser) ToUser(userId uint64) User {
	// unmarshal avatar
	avatar := Avatar{}
	_ = avatar.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, u.Avatar))) // this is quite hacky

	return User{
		Id:          userId,
		Username:    u.Username,
		GlobalName:  u.GlobalName,
		Avatar:      avatar,
		Bot:         u.Bot,
		Flags:       u.Flags,
		PremiumType: u.PremiumType,
	}
}
