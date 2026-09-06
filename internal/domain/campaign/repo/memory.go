// Package repo provides an in-memory campaign.Repo (tests + local dev).
package repo

import (
	"context"
	"sync"

	"github.com/tatayooo/review-poc-go/internal/domain/campaign"
)

// Memory is a concurrency-safe in-memory campaign store.
type Memory struct {
	mu   sync.RWMutex
	byID map[string]*campaign.Campaign
}

// NewMemory constructs an empty store.
func NewMemory() *Memory { return &Memory{byID: map[string]*campaign.Campaign{}} }

// Get returns a deep copy to keep the aggregate immutable outside the store.
func (m *Memory) Get(_ context.Context, id string) (*campaign.Campaign, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.byID[id]
	if !ok {
		return nil, campaign.ErrNotFound
	}
	cp := *c
	return &cp, nil
}

// List returns all campaigns for an advertiser, oldest first.
func (m *Memory) List(_ context.Context, advertiserID string) ([]*campaign.Campaign, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*campaign.Campaign, 0, len(m.byID))
	for _, c := range m.byID {
		if c.AdvertiserID == advertiserID {
			cp := *c
			out = append(out, &cp)
		}
	}
	return out, nil
}

// Save inserts or updates by ID.
func (m *Memory) Save(_ context.Context, c *campaign.Campaign) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *c
	m.byID[c.ID] = &cp
	return nil
}
