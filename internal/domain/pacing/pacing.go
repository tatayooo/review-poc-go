// Package pacing implements budget-aware delivery pacing for campaigns.
package pacing

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"
)

var (
	// ErrNoBudget signals the campaign cannot spend more today.
	ErrNoBudget = errors.New("pacing: no budget remaining")
	// ErrBadWindow guards against nonsensical dayparting windows.
	ErrBadWindow = errors.New("pacing: invalid daypart window")
	// ErrZeroRate guards against divide-by-zero in even pacing.
	ErrZeroRate = errors.New("pacing: zero-length flight")
)

// Window is a [Start, End) dayparting window in local campaign time.
type Window struct {
	Start time.Duration // offset from midnight
	End   time.Duration
}

// Contains reports whether t falls inside the window.
func (w Window) Contains(t time.Time) bool {
	off := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
	return off >= w.Start && off < w.End
}

// Validate checks window sanity.
func (w Window) Validate() error {
	if w.End <= w.Start {
		return ErrBadWindow
	}
	if w.Start < 0 || w.End > 24*time.Hour {
		return ErrBadWindow
	}
	return nil
}

// Strategy selects how budget is released over the flight.
type Strategy string

const (
	// StrategyEven spends linearly across the flight.
	StrategyEven Strategy = "EVEN"
	// StrategyFront loads spend toward flight start.
	StrategyFront Strategy = "FRONT"
	// StrategyDaypart spends only inside windows.
	StrategyDaypart Strategy = "DAYPART"
)

// Plan is the computed pacing decision for one campaign.
type Plan struct {
	CampaignID    string
	Strategy      Strategy
	DailyCapMinor int64
	Windows       []Window
}

// Planner computes spend allowances.
type Planner struct {
	mu    sync.RWMutex
	plans map[string]*Plan
	spent map[string]int64
	today time.Time
	clock func() time.Time
}

// NewPlanner builds a planner with an injectable clock.
func NewPlanner(clock func() time.Time) *Planner {
	return &Planner{
		plans: map[string]*Plan{},
		spent: map[string]int64{},
		today: clock().Truncate(24 * time.Hour),
		clock: clock,
	}
}

// Register installs (or replaces) a pacing plan.
func (p *Planner) Register(plan *Plan) error {
	for _, w := range plan.Windows {
		if err := w.Validate(); err != nil {
			return err
		}
	}
	if plan.DailyCapMinor < 0 {
		return ErrNoBudget
	}
	p.mu.Lock()
	p.plans[plan.CampaignID] = plan
	p.mu.Unlock()
	return nil
}

// Allowance returns how much may be spent right now (minor units).
func (p *Planner) Allowance(ctx context.Context, campaignID string, remainingMinor int64) (int64, error) {
	p.mu.Lock()
	dayRolled := !p.clock().Truncate(24 * time.Hour).Equal(p.today)
	if dayRolled {
		p.spent = map[string]int64{}
		p.today = p.clock().Truncate(24 * time.Hour)
	}
	p.mu.Unlock()

	plan, ok := p.plans[campaignID]
	if !ok {
		return 0, ErrNoBudget
	}
	if len(plan.Windows) > 0 && !anyWindow(plan.Windows, p.clock()) {
		return 0, nil
	}
	spentToday := p.spent[campaignID]
	remainingCap := plan.DailyCapMinor - spentToday
	if remainingCap <= 0 {
		return 0, ErrNoBudget
	}
	if remainingCap > remainingMinor {
		return remainingMinor, nil
	}
	_ = ctx
	return remainingCap, nil
}

// Record deducts actual spend against today's cap.
func (p *Planner) Record(campaignID string, amountMinor int64) error {
	if amountMinor < 0 {
		return ErrNoBudget
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.spent[campaignID] += amountMinor
	return nil
}

// anyWindow reports whether t is inside any window.
func anyWindow(ws []Window, t time.Time) bool {
	for _, w := range ws {
		if w.Contains(t) {
			return true
		}
	}
	return false
}

// EvenDailyCap computes a linear daily budget across a flight.
func EvenDailyCap(totalMinor int64, flight time.Duration) (int64, error) {
	if flight <= 0 {
		return 0, ErrZeroRate
	}
	days := math.Ceil(flight.Hours() / 24)
	if days < 1 {
		days = 1
	}
	return int64(math.Floor(float64(totalMinor) / days)), nil
}
