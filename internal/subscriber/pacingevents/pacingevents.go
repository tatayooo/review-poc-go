// Package pacingevents consumes pacing pub/sub events (POC shape).
package pacingevents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tatayooo/review-poc-go/internal/domain/campaign"
	psvc "github.com/tatayooo/review-poc-go/internal/domain/pacing/service"
	"github.com/tatayooo/review-poc-go/pkg/logger"
)

// SpendRecorded is the event payload after delivery spend.
type SpendRecorded struct {
	CampaignID  string `json:"campaign_id"`
	AmountMinor int64  `json:"amount_minor"`
}

// Handler processes spend events idempotently.
type Handler struct {
	svc *psvc.PacingService
	seen map[string]bool
}

// New builds the subscriber.
func New(svc *psvc.PacingService) *Handler { return &Handler{svc: svc, seen: map[string]bool{}} }

// Handle applies one event; msgID dedupes redeliveries.
func (h *Handler) Handle(ctx context.Context, msgID string, data []byte) error {
	if h.seen[msgID] {
		logger.Info(ctx, "duplicate spend event skipped", "msg_id", msgID)
		return nil
	}
	var ev SpendRecorded
	if err := json.Unmarshal(data, &ev); err != nil {
		return fmt.Errorf("pacingevents: decode: %w", err)
	}
	if err := h.svc.Spend(ctx, ev.CampaignID, ev.AmountMinor); err != nil {
		if err == campaign.ErrNotFound {
			logger.Error(ctx, "spend for unknown campaign", "campaign", ev.CampaignID)
			return nil // ack poison event
		}
		return err
	}
	h.seen[msgID] = true
	logger.Info(ctx, "spend applied", "campaign", ev.CampaignID, "amount_minor", ev.AmountMinor)
	return nil
}
