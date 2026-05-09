package search

import (
	"strings"

	"github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"
)

type IndexedMessage struct {
	MessageID      int64
	ConversationID int64
	Tokens         []string
	Status         string
}

func BuildIndexRecord(msg model.Message) IndexedMessage {
	return IndexedMessage{MessageID: msg.ID, ConversationID: msg.ConversationID, Tokens: tokenize(msg.Content), Status: msg.Status}
}

func tokenize(s string) []string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return nil
	}
	return strings.Fields(s)
}

func MightBeDuplicateByTokens(a, b IndexedMessage) bool {
	if a.ConversationID != b.ConversationID {
		return false
	}
	if len(a.Tokens) == 0 || len(b.Tokens) == 0 {
		return false
	}
	matches := 0
	seen := map[string]bool{}
	for _, t := range a.Tokens {
		seen[t] = true
	}
	for _, t := range b.Tokens {
		if seen[t] {
			matches++
		}
	}
	return matches >= len(a.Tokens)/2+1
}
