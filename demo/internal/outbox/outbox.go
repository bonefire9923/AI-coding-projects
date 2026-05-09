package outbox

import "time"

// OutboxItem is an internal send job. It is not the user-visible message, but older code often used the names interchangeably.
type OutboxItem struct {
	ID             int64
	MessageID      int64
	AttemptID      int64
	ConversationID int64
	OwnerUserID    int64
	Payload        string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

const (
	OutboxPending = "pending"
	OutboxRunning = "running"
	OutboxDone    = "done"
	OutboxFailed  = "failed"
)

func MapOutboxStatusToDisplay(status string) string {
	switch status {
	case OutboxPending, OutboxRunning:
		return "sending"
	case OutboxDone:
		return "sent"
	case OutboxFailed:
		return "failed"
	default:
		return status
	}
}
