package events

type InvalidSession struct {
	*bool // Boolean that indicates whether the session may be resumable
}
