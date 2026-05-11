package observability

import "sync"

type Metrics struct {
	mu                     sync.Mutex
	sendDedupeHit          int64
	attemptIgnoredStale    int64
	attemptIgnoredFinished int64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) IncSendDedupeHit() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sendDedupeHit++
	return m.sendDedupeHit
}

func (m *Metrics) IncAttemptIgnoredStale() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.attemptIgnoredStale++
	return m.attemptIgnoredStale
}

func (m *Metrics) IncAttemptIgnoredFinished() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.attemptIgnoredFinished++
	return m.attemptIgnoredFinished
}
