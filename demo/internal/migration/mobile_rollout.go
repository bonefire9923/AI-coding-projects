package migration

import "time"

type RolloutStep struct {
	Name      string
	StartedAt time.Time
	Percent   int
	Notes     string
}

func DefaultRollout() []RolloutStep {
	return []RolloutStep{
		{Name: "legacy-status-display", StartedAt: time.Now().Add(-72 * time.Hour), Percent: 100, Notes: "list display reads compatibility status"},
		{Name: "retry-as-new-item", StartedAt: time.Now().Add(-48 * time.Hour), Percent: 60, Notes: "keeps old retry UI behavior"},
		{Name: "server-sync-cursor", StartedAt: time.Now().Add(-24 * time.Hour), Percent: 50, Notes: "server stores last returned cursor"},
	}
}

func IsStepEnabled(steps []RolloutStep, name string) bool {
	for _, step := range steps {
		if step.Name == name && step.Percent > 0 {
			return true
		}
	}
	return false
}
