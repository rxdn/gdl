package user

import "fmt"

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

func (u *User) AvatarUrl(size int) string {
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%d/%s.webp?size=%d", u.Id, u.Avatar, size)
}

func (u *User) Mention() string {
	return fmt.Sprintf("<@%d>", u.Id)
}

// converts a discrim of 1 => 0001
func (u *User) PadDiscriminator() string {
	return fmt.Sprintf("%04d", u.Discriminator)
}

func (u *User) ToCachedUser() CachedUser {
	return CachedUser{
		Username:      u.Username,
		Discriminator: u.Discriminator,
		Avatar:        u.Avatar,
		Bot:           u.Bot,
		Flags:         u.Flags,
		PremiumType:   u.PremiumType,
	}
}