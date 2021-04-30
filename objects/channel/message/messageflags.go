package message

type MessageFlag uint

const (
	FlagCrossposted MessageFlag = 1 << iota
	FlagIsCrosspost
	FlagSupressEmbeds
	FlagSourceMessageDeleted
	FlagUrgent
	_ // 1 << 5 not documented
	FlagEphemeral
	FlagLoading
)

func SumFlags(flags ...MessageFlag) (sum uint) {
	for _, flag := range flags {
		sum += uint(flag)
	}

	return
}
