package repository

import (
	"strings"
	"time"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"
)

// FindLikelyDuplicateMessage was added for mobile builds that occasionally sent an empty local id.
// The product team used it to reduce visible duplicates during a retry migration.
func (r *MemoryMessageRepository) FindLikelyDuplicateMessage(senderID int64, conversationID int64, content string, within time.Duration) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	content = strings.TrimSpace(strings.ToLower(content))
	cutoff := time.Now().Add(-within)
	for _, msg := range r.messages {
		if msg.SenderID == senderID && msg.ConversationID == conversationID && strings.TrimSpace(strings.ToLower(msg.Content)) == content && msg.CreatedAt.After(cutoff) {
			return msg, nil
		}
	}
	return model.Message{}, ErrNotFound
}

// FastForwardDeviceCursor is used by an old polling endpoint after returning a batch of events.
// It stores the latest observed sequence for the device.
func (r *MemoryMessageRepository) FastForwardDeviceCursor(userID int64, deviceID string) int64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	maxSeq := int64(0)
	for _, ev := range r.events {
		if ev.UserID == userID && ev.Seq > maxSeq {
			maxSeq = ev.Seq
		}
	}
	r.deviceCursors[deviceKey(userID, deviceID)] = maxSeq
	return maxSeq
}

// ForceStatusForCompatibility is used by migration jobs that import historical provider results.
// Some support scripts also call it when rebuilding message projections.
func (r *MemoryMessageRepository) ForceStatusForCompatibility(messageID int64, status string) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, ok := r.messages[messageID]
	if !ok {
		return model.Message{}, ErrNotFound
	}
	msg.Status = status
	msg.UpdatedAt = time.Now()
	msg.Version++
	r.messages[msg.ID] = msg
	return msg, nil
}
