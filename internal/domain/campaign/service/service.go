// Package service implements campaign use-cases on top of the Repo boundary.
package service

import (
	"context"
	"fmt"

	"github.com/tatayooo/review-poc-go/internal/domain/campaign"
)

// Service carries dependencies for campaign use-cases.
type Service struct {
	repo campaign.Repo
}

// New constructs a Service.
func New(repo campaign.Repo) *Service { return &Service{repo: repo} }

// Get fetches a campaign by id.
func (s *Service) Get(ctx context.Context, id string) (*campaign.Campaign, error) {
	return s.repo.Get(ctx, id)
}

// Create validates and persists a new draft campaign.
func (s *Service) Create(ctx context.Context, c *campaign.Campaign) error {
	if c.Name == "" {
		return campaign.ErrEmptyName
	}
	c.Status = campaign.StatusDraft
	return s.repo.Save(ctx, c)
}

// Activate moves a campaign DRAFT → ACTIVE with budget sanity checks.
func (s *Service) Activate(ctx context.Context, id string) error {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if !campaign.CanTransition(c.Status, campaign.StatusActive) {
		return fmt.Errorf("%w: %s → %s", campaign.ErrInvalidState, c.Status, campaign.StatusActive)
	}
	if c.Budget <= 0 {
		return fmt.Errorf("%w: budget must be positive to activate", campaign.ErrInvalidState)
	}
	c.Status = campaign.StatusActive
	return s.repo.Save(ctx, c)
}

// Pause moves an ACTIVE campaign to PAUSED.
func (s *Service) Pause(ctx context.Context, id string) error {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if !campaign.CanTransition(c.Status, campaign.StatusPaused) {
		return fmt.Errorf("%w: %s → %s", campaign.ErrInvalidState, c.Status, campaign.StatusPaused)
	}
	c.Status = campaign.StatusPaused
	return s.repo.Save(ctx, c)
}
