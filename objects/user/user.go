package user

import "fmt"

type User struct {
	Id            uint64 `json:"id,string"`
	Username      string `json:"username"`
	Discriminator uint16 `json:"discriminator,string"`
	Avatar        Avatar `json:"avatar"`
	Bot           bool   `json:"bot"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         uint32 `json:"flags"`
	PremiumType   int    `json:"premium_type"`
}

// shortcut, ignores errors
func (u *User) AvatarUrl(size int) string {
	avatar, _ := u.Avatar.String()

	var extension string
	if u.Avatar.Animated {
		extension = "gif"
	} else {
		extension = "webp"
	}

	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%d/%s.%s?size=%d", u.Id, avatar, extension, size)
}

func (u *User) Mention() string {
	return fmt.Sprintf("<@%d>", u.Id)
}

// converts a discrim of 1 => 0001
func (u *User) PadDiscriminator() string {
	return fmt.Sprintf("%04d", u.Discriminator)
}

func (u *User) ToCachedUser() CachedUser {
	avatar, _ := u.Avatar.String()

	return CachedUser{
		Username:      u.Username,
		Discriminator: u.Discriminator,
		Avatar:        avatar,
		Bot:           u.Bot,
		Flags:         u.Flags,
		PremiumType:   u.PremiumType,
	}
}
