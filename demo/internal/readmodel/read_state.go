package readmodel

import "time"

// ReadState belongs to a user in a conversation. It should not be confused with message delivery.
type ReadState struct {
	UserID             int64
	ConversationID     int64
	LastReadMessageID  int64
	LastReadEventSeq   int64
	UpdatedAt          time.Time
}

// ReadStateProjection stores only what the recent-chat list needs.
type ReadStateProjection struct {
	UserID         int64
	ConversationID int64
	UnreadCount    int
	UpdatedAt      time.Time
}

func MergeReadProjection(old ReadStateProjection, delta int) ReadStateProjection {
	old.UnreadCount += delta
	if old.UnreadCount < 0 {
		old.UnreadCount = 0
	}
	old.UpdatedAt = time.Now()
	return old
}
