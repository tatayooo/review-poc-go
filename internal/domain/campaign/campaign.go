// Package campaign holds the campaign domain: entities, repo interface,
// service errors. Pure Go — no gRPC, no DB imports (DDD layering).
package campaign

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotFound is returned when a campaign does not exist.
	ErrNotFound = errors.New("campaign not found")
	// ErrInvalidState is returned on illegal lifecycle transitions.
	ErrInvalidState = errors.New("invalid campaign state")
	// ErrEmptyName guards against unnamed campaigns.
	ErrEmptyName = errors.New("campaign name must not be empty")
)

// Status is the campaign lifecycle state.
type Status string

const (
	StatusDraft     Status = "DRAFT"
	StatusActive    Status = "ACTIVE"
	StatusPaused    Status = "PAUSED"
	StatusArchived  Status = "ARCHIVED"
)

// Campaign is the aggregate root.
type Campaign struct {
	ID        string
	Name      string
	AdvertiserID string
	Status    Status
	Budget    int64 // minor units
	Spent     int64 // minor units
	StartAt   time.Time
	EndAt     time.Time
}

// CanTransition reports whether from → to is a legal lifecycle move.
func CanTransition(from, to Status) bool {
	switch from {
	case StatusDraft:
		return to == StatusActive
	case StatusActive:
		return to == StatusPaused || to == StatusArchived
	case StatusPaused:
		return to == StatusActive || to == StatusArchived
	default:
		return false
	}
}

// Remaining reports unspent budget in minor units.
func (c *Campaign) Remaining() int64 { return c.Budget - c.Spent }

// Repo is the persistence boundary implemented by infra adapters.
type Repo interface {
	Get(ctx context.Context, id string) (*Campaign, error)
	List(ctx context.Context, advertiserID string) ([]*Campaign, error)
	Save(ctx context.Context, c *Campaign) error
}
