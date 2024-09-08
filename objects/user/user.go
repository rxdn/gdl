package user

import "fmt"

type User struct {
	Id          uint64  `json:"id,string"`
	Username    string  `json:"username"`
	GlobalName  *string `json:"global_name"`
	Avatar      Avatar  `json:"avatar"`
	Bot         bool    `json:"bot"`
	MfaEnabled  bool    `json:"mfa_enabled"`
	Locale      string  `json:"locale"`
	Verified    bool    `json:"verified"`
	Email       string  `json:"email"`
	Flags       uint32  `json:"flags"`
	PremiumType int     `json:"premium_type"`
}

// shortcut, ignores errors
func (u *User) AvatarUrl(size int) string {
	hash := u.Avatar.String()
	// if blank avatar, return a blank string so that we can use omitempty
	if len(hash) == 0 {
		return ""
	}

	var extension string
	if u.Avatar.Animated {
		extension = "gif"
	} else {
		extension = "webp"
	}

	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%d/%s.%s?size=%d", u.Id, u.Avatar.String(), extension, size)
}

func (u *User) EffectiveName() string {
	if u.GlobalName != nil && len(*u.GlobalName) > 0 {
		return *u.GlobalName
	}

	return u.Username
}

func (u *User) Mention() string {
	return fmt.Sprintf("<@%d>", u.Id)
}

func (u *User) ToCachedUser() CachedUser {
	return CachedUser{
		Username:    u.Username,
		GlobalName:  u.GlobalName,
		Avatar:      u.Avatar.String(),
		Bot:         u.Bot,
		Flags:       u.Flags,
		PremiumType: u.PremiumType,
	}
}
