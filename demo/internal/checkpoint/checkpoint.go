package checkpoint

import "time"

// Checkpoint is shared by several old jobs. Some jobs use it for sync, some for read state, some for export jobs.
type Checkpoint struct {
	OwnerUserID    int64
	DeviceID       string
	ConversationID int64
	Kind           string
	Value          int64
	UpdatedAt      time.Time
}

const (
	KindSyncReturned = "sync_returned"
	KindSyncApplied  = "sync_applied"
	KindReadMessage  = "read_message"
	KindExported     = "exported"
)

func ShouldReuseCheckpoint(a Checkpoint, b Checkpoint) bool {
	return a.OwnerUserID == b.OwnerUserID && a.DeviceID == b.DeviceID && a.ConversationID == b.ConversationID
}
