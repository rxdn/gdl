package guild

type ExplicitContentFilterLevel int

const (
	DISABLED              ExplicitContentFilterLevel = 0
	MEMBERS_WITHOUT_ROLES ExplicitContentFilterLevel = 1
	ALL_MEMBERS           ExplicitContentFilterLevel = 2
)
