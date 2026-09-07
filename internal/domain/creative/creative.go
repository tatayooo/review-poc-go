// Package creative manages ad creatives and their review lifecycle.
package creative

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrRejected means the creative failed review.
	ErrRejected = errors.New("creative: rejected")
	// ErrPendingReview means no decision yet.
	ErrPendingReview = errors.New("creative: pending review")
	// ErrBadAsset means the asset failed validation.
	ErrBadAsset = errors.New("creative: invalid asset")
)

var safeURLRe = regexp.MustCompile(`^https://[a-z0-9.-]+/.*$`)

// Format enumerates asset kinds.
type Format string

const (
	FormatImage Format = "IMAGE"
	FormatVideo Format = "VIDEO"
	FormatHTML  Format = "HTML"
)

// Aspect groups sizes for responsive slots.
type Aspect struct {
	W int
	H int
}

// Asset is one renderable unit of a creative.
type Asset struct {
	URL      string
	Format   Format
	Aspect   Aspect
	BytesMax int64
}

// Validate checks scheme, format and byte cap.
func (a *Asset) Validate() error {
	if !safeURLRe.MatchString(a.URL) {
		return fmt.Errorf("%w: url must be https and absolute", ErrBadAsset)
	}
	switch a.Format {
	case FormatImage, FormatVideo, FormatHTML:
	default:
		return fmt.Errorf("%w: unknown format %q", ErrBadAsset, a.Format)
	}
	if a.BytesMax <= 0 {
		return fmt.Errorf("%w: byte cap must be positive", ErrBadAsset)
	}
	if a.Aspect.W <= 0 || a.Aspect.H <= 0 {
		return fmt.Errorf("%w: aspect must be positive", ErrBadAsset)
	}
	return nil
}

// ReviewState is the moderation lifecycle.
type ReviewState string

const (
	ReviewPending  ReviewState = "PENDING"
	ReviewApproved ReviewState = "APPROVED"
	ReviewRejected ReviewState = "REJECTED"
)

// Creative is the aggregate: assets plus review state plus targeting hints.
type Creative struct {
	ID          string
	CampaignID  string
	Title       string
	Assets      []Asset
	State       ReviewState
	RejectCode  string
	LandingURL  string
	UpdatedAt   time.Time
}

// Validate runs all invariants before persistence.
func (c *Creative) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("%w: title required", ErrBadAsset)
	}
	if !safeURLRe.MatchString(c.LandingURL) {
		return fmt.Errorf("%w: landing url must be https", ErrBadAsset)
	}
	if len(c.Assets) == 0 {
		return fmt.Errorf("%w: at least one asset required", ErrBadAsset)
	}
	for i := range c.Assets {
		if err := c.Assets[i].Validate(); err != nil {
			return fmt.Errorf("asset %d: %w", i, err)
		}
	}
	return nil
}

// Approve flips state with transition guard.
func (c *Creative) Approve() error {
	if c.State != ReviewPending {
		return fmt.Errorf("creative: approve from %s: %w", c.State, ErrRejected)
	}
	c.State = ReviewApproved
	c.RejectCode = ""
	return nil
}

// Reject records a moderation refusal.
func (c *Creative) Reject(code string) error {
	if c.State != ReviewPending {
		return fmt.Errorf("creative: reject from %s: %w", c.State, ErrRejected)
	}
	if code == "" {
		return fmt.Errorf("%w: reject code required", ErrRejected)
	}
	c.State = ReviewRejected
	c.RejectCode = code
	return nil
}

// Renderable reports whether the creative may serve.
func (c *Creative) Renderable() bool { return c.State == ReviewApproved }

// Repo is the persistence boundary.
type Repo interface {
	Get(ctx context.Context, id string) (*Creative, error)
	ListByCampaign(ctx context.Context, campaignID string) ([]*Creative, error)
	Save(ctx context.Context, c *Creative) error
}
