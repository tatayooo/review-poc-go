// Package service implements creative moderation use-cases.
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/tatayooo/review-poc-go/internal/domain/creative"
)

// bannedWords is a stub moderation list (real impl: policy service).
var bannedWords = []string{"free money", "guaranteed win", "click here now"}

// ModerationService validates and moderates creatives.
type ModerationService struct {
	repo creative.Repo
}

// NewModeration wires the service.
func NewModeration(repo creative.Repo) *ModerationService { return &ModerationService{repo: repo} }

// Submit validates and stores a creative in PENDING.
func (s *ModerationService) Submit(ctx context.Context, c *creative.Creative) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if err := scanBanned(c.Title); err != nil {
		return err
	}
	c.State = creative.ReviewPending
	return s.repo.Save(ctx, c)
}

// AutoModerate applies deterministic rules; approves or rejects.
func (s *ModerationService) AutoModerate(ctx context.Context, id string) error {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := scanBanned(c.Title); err != nil {
		_ = c.Reject("BANNED_PHRASE")
		return s.repo.Save(ctx, c)
	}
	for _, a := range c.Assets {
		if a.Format == creative.FormatHTML && strings.Contains(a.URL, "doubleclick") {
			_ = c.Reject("HTML_VENDOR")
			return s.repo.Save(ctx, c)
		}
	}
	return c.Approve()
}

// Serveable lists approved creatives for a campaign.
func (s *ModerationService) Serveable(ctx context.Context, campaignID string) ([]*creative.Creative, error) {
	all, err := s.repo.ListByCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	out := make([]*creative.Creative, 0, len(all))
	for _, c := range all {
		if c.Renderable() {
			out = append(out, c)
		}
	}
	return out, nil
}

// scanBanned enforces the wordlist.
func scanBanned(title string) error {
	lower := strings.ToLower(title)
	for _, w := range bannedWords {
		if strings.Contains(lower, w) {
			return fmt.Errorf("creative: banned phrase %q: %w", w, creative.ErrRejected)
		}
	}
	return nil
}
