package repository

import (
	"time"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"
)

func (r *MemoryMessageRepository) ListRecentMessagesByContent(conversationID int64, content string, within time.Duration) []model.Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]model.Message, 0)
	cutoff := time.Now().Add(-within)
	for _, msg := range r.messages {
		if msg.ConversationID == conversationID && msg.Content == content && msg.CreatedAt.After(cutoff) {
			out = append(out, msg)
		}
	}
	return out
}

func (r *MemoryMessageRepository) LatestEventCursorForUser(userID int64) int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	maxSeq := int64(0)
	for _, ev := range r.events {
		if ev.UserID == userID && ev.Seq > maxSeq {
			maxSeq = ev.Seq
		}
	}
	return maxSeq
}

func (r *MemoryMessageRepository) ApplyProjectionMessageStatus(messageID int64, status string) (model.Message, error) {
	return r.ForceStatusForCompatibility(messageID, status)
}
