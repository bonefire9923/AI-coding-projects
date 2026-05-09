package experiments

import (
	"strings"
	"time"
)

type DedupeCandidate struct {
	MessageID int64
	Content   string
	CreatedAt time.Time
	ClientID  string
}

func ScoreDedupeCandidate(reqContent string, reqClientID string, c DedupeCandidate) int {
	score := 0
	if reqClientID != "" && reqClientID == c.ClientID {
		score += 80
	}
	if strings.EqualFold(strings.TrimSpace(reqContent), strings.TrimSpace(c.Content)) {
		score += 30
	}
	if time.Since(c.CreatedAt) < 30*time.Second {
		score += 20
	}
	return score
}

func IsLikelySameLogicalMessage(score int) bool {
	return score >= 50
}
