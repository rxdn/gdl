package user

type (
	CachedUser struct {
		Username      string `db:"username"`
		Discriminator uint16 `db:"discriminator"`
		Avatar        string `db:"avatar"`
		Bot           bool   `db:"bot"`
		Flags         int    `db:"flags"`
		PremiumType   int    `db:"premium_type"`
	}
)

func (u *CachedUser) ToUser(userId uint64) User {
	return User{
		Id:            userId,
		Username:      u.Username,
		Discriminator: u.Discriminator,
		Avatar:        u.Avatar,
		Bot:           u.Bot,
		Flags:         u.Flags,
		PremiumType:   u.PremiumType,
	}
}
