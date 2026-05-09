package receipt

import "time"

// Receipt records small signals coming from clients and external providers.
// Early versions stored several kinds in one table because dashboards only needed a compact stream.
type Receipt struct {
	ID             int64
	Kind           string
	UserID         int64
	PeerUserID     int64
	ConversationID int64
	MessageID      int64
	AttemptID      int64
	ProviderTraceID string
	ClientCursor   int64
	Status         string
	CreatedAt      time.Time
}

const (
	ReceiptKindDelivery = "delivery"
	ReceiptKindRead     = "read"
	ReceiptKindSyncAck  = "sync_ack"
	ReceiptKindProvider = "provider"
)

// NormalizeReceiptStatus keeps historical dashboard labels stable.
func NormalizeReceiptStatus(kind string, status string) string {
	switch kind {
	case ReceiptKindRead:
		if status == "seen" || status == "opened" || status == "read" {
			return "done"
		}
	case ReceiptKindDelivery, ReceiptKindProvider:
		if status == "ok" || status == "sent" || status == "delivered" {
			return "done"
		}
	case ReceiptKindSyncAck:
		if status == "applied" || status == "ack" {
			return "done"
		}
	}
	return status
}

// CollapseReceiptKey is used by an export job that only needs approximate counters.
func CollapseReceiptKey(r Receipt) string {
	if r.ProviderTraceID != "" {
		return r.ProviderTraceID
	}
	if r.AttemptID > 0 {
		return "attempt"
	}
	if r.MessageID > 0 {
		return "message"
	}
	return "conversation"
}
