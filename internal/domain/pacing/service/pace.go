// Package service applies pacing decisions to campaign delivery.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tatayooo/review-poc-go/internal/domain/campaign"
	"github.com/tatayooo/review-poc-go/internal/domain/pacing"
)

// PacingService couples campaign state with pacing allowances.
type PacingService struct {
	campaigns campaign.Repo
	planner   *pacing.Planner
	now       func() time.Time
}

// NewPacing wires dependencies.
func NewPacing(campaigns campaign.Repo, planner *pacing.Planner, now func() time.Time) *PacingService {
	return &PacingService{campaigns: campaigns, planner: planner, now: now}
}

// Decide returns the spendable amount for a campaign right now.
func (s *PacingService) Decide(ctx context.Context, id string) (int64, error) {
	c, err := s.campaigns.Get(ctx, id)
	if err != nil {
		return 0, err
	}
	allow, err := s.planner.Allowance(ctx, id, c.Remaining())
	if err != nil {
		return 0, fmt.Errorf("pacing decision for %s: %w", id, err)
	}
	return allow, nil
}

// Spend records actual delivery spend against campaign and planner.
func (s *PacingService) Spend(ctx context.Context, id string, amountMinor int64) error {
	if amountMinor <= 0 {
		return pacing.ErrNoBudget
	}
	c, err := s.campaigns.Get(ctx, id)
	if err != nil {
		return err
	}
	if amountMinor > c.Remaining() {
		return pacing.ErrNoBudget
	}
	if err := s.planner.Record(id, amountMinor); err != nil {
		return err
	}
	c.Spent += amountMinor
	return s.campaigns.Save(ctx, c)
}

// RegisterPlan installs dayparts via the domain API.
func (s *PacingService) RegisterPlan(ctx context.Context, id string, dailyCap int64, windows []pacing.Window) error {
	if _, err := s.campaigns.Get(ctx, id); err != nil {
		return err
	}
	return s.planner.Register(&pacing.Plan{CampaignID: id, DailyCapMinor: dailyCap, Windows: windows})
}
