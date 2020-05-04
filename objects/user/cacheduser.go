package user

import "fmt"

type (
	CachedUser struct {
		Username      string `db:"username"`
		Discriminator uint16 `db:"discriminator"`
		Avatar        string `db:"avatar"`
		Bot           bool   `db:"bot"`
		Flags         uint32 `db:"flags"`
		PremiumType   int    `db:"premium_type"`
	}
)

func (u *CachedUser) ToUser(userId uint64) User {
	// unmarshal avatar
	avatar := Avatar{}
	_ = avatar.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, u.Avatar))) // this is quite hacky

	return User{
		Id:            userId,
		Username:      u.Username,
		Discriminator: u.Discriminator,
		Avatar:        avatar,
		Bot:           u.Bot,
		Flags:         u.Flags,
		PremiumType:   u.PremiumType,
	}
}
