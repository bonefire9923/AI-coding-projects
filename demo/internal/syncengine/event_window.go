package syncengine

import "github.com/example/backend-ai-coding-challenge-demo-v6/internal/model"

type Window struct {
	Events     []model.SyncEvent
	NextCursor int64
	HasMore    bool
}

func BuildWindow(events []model.SyncEvent, limit int) Window {
	if limit <= 0 || limit > len(events) {
		limit = len(events)
	}
	out := Window{Events: events[:limit]}
	for _, ev := range out.Events {
		if ev.Seq > out.NextCursor {
			out.NextCursor = ev.Seq
		}
	}
	out.HasMore = len(events) > limit
	return out
}

func CompactByMessage(events []model.SyncEvent) []model.SyncEvent {
	latest := make(map[int64]model.SyncEvent)
	for _, ev := range events {
		if old, ok := latest[ev.MessageID]; !ok || ev.Seq > old.Seq {
			latest[ev.MessageID] = ev
		}
	}
	out := make([]model.SyncEvent, 0, len(latest))
	for _, ev := range latest {
		out = append(out, ev)
	}
	return out
}
