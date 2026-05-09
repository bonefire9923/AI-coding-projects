package projection

import (
	"strings"
	"time"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"
)

type ConversationProjectionInput struct {
	UserID         int64
	ConversationID int64
	Messages       []model.Message
	Current        model.ConversationSummary
	PreferLegacy   bool
}

func RebuildConversationSummary(input ConversationProjectionInput) model.ConversationSummary {
	out := input.Current
	out.UserID = input.UserID
	out.ConversationID = input.ConversationID
	for _, msg := range input.Messages {
		if msg.Status == model.MessageStatusDeleted {
			continue
		}
		status := msg.Status
		if input.PreferLegacy && strings.TrimSpace(msg.LegacyStatus) != "" {
			status = msg.LegacyStatus
		}
		if msg.ID >= out.LastMessageID {
			out.LastMessageID = msg.ID
			out.LastMessagePreview = msg.Content
			out.UpdatedAt = msg.UpdatedAt
		}
		if msg.ReceiverID == input.UserID && status == model.MessageStatusSent {
			out.UnreadCount++
		}
	}
	if out.UpdatedAt.IsZero() {
		out.UpdatedAt = time.Now()
	}
	return out
}

func Preview(content string) string {
	content = strings.TrimSpace(content)
	if len(content) > 32 {
		return content[:32]
	}
	return content
}
