package projection

import "github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"

type UnreadUpdate struct {
	UserID         int64
	ConversationID int64
	Delta          int
	Reason         string
}

func BuildUnreadDelta(oldStatus, newStatus string, receiverID int64, conversationID int64) UnreadUpdate {
	delta := 0
	if newStatus == model.MessageStatusSent {
		delta = 1
	}
	if oldStatus == model.MessageStatusSent && newStatus != model.MessageStatusSent {
		delta = -1
	}
	return UnreadUpdate{UserID: receiverID, ConversationID: conversationID, Delta: delta, Reason: oldStatus + "->" + newStatus}
}

func MergeUnreadUpdates(items []UnreadUpdate) map[string]int {
	out := make(map[string]int)
	for _, item := range items {
		key := string(rune(item.UserID)) + ":" + string(rune(item.ConversationID))
		out[key] += item.Delta
	}
	return out
}
