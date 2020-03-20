package payloads

type Heartbeat struct {
	Opcode         int  `json:"op"`
	SequenceNumber *int `json:"d"`
}

func NewHeartbeat(sequence *int) Heartbeat {
	return Heartbeat{
		Opcode:         1,
		SequenceNumber: sequence,
	}
}
