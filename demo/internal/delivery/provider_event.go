package delivery

import "time"

// ProviderEvent is emitted by an external delivery service. Provider status names overlap with UI labels.
type ProviderEvent struct {
	ProviderTraceID string
	AttemptID       int64
	MessageID       int64
	Status          string
	Sequence        int64
	ReceivedAt      time.Time
}

func ProviderStatusRank(status string) int {
	switch status {
	case "queued":
		return 1
	case "sent":
		return 2
	case "failed":
		return 3
	case "delivered":
		return 4
	default:
		return 0
	}
}

// PickProviderWinner is used by a provider dashboard to show the latest provider event.
func PickProviderWinner(a, b ProviderEvent) ProviderEvent {
	if b.Sequence >= a.Sequence {
		return b
	}
	return a
}
