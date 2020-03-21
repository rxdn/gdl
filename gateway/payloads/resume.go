package payloads

type (
	Resume struct {
		Opcode int        `json:"op"`
		Data   ResumeData `json:"d"`
	}

	ResumeData struct {
		Token          string `json:"token"`
		SessionId      string `json:"session_id"`
		SequenceNumber int    `json:"seq"`
	}
)

func NewResume(token, sessionId string, sequenceNumber int) Resume {
	payload := Resume{
		Opcode: 6,
		Data: ResumeData{
			Token:          token,
			SessionId:      sessionId,
			SequenceNumber: sequenceNumber,
		},
	}
	return payload
}
