// Package repo keeps pacing telemetry (in-memory for the POC).
package repo

import "sync"

// Metrics is a cheap counter registry keyed by campaign.
type Metrics struct {
	mu     sync.Mutex
	counts map[string]map[string]int64
}

// NewMetrics builds an empty registry.
func NewMetrics() *Metrics { return &Metrics{counts: map[string]map[string]int64{}} }

// Inc bumps counter for campaign.
func (m *Metrics) Inc(campaign, name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.counts[campaign] == nil {
		m.counts[campaign] = map[string]int64{}
	}
	m.counts[campaign][name]++
}

// Snapshot copies all counters.
func (m *Metrics) Snapshot() map[string]map[string]int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[string]map[string]int64, len(m.counts))
	for k, v := range m.counts {
		out[k] = map[string]int64{}
		for ck, cv := range v {
			out[k][ck] = cv
		}
	}
	return out
}
