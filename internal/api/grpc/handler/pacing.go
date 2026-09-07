// Pacing transport endpoints.
package handler

import (
	"context"
	"time"

	campaignv1 "github.com/tatayooo/review-poc-go/internal/api/grpc/campaign/v1"
	"github.com/tatayooo/review-poc-go/internal/domain/pacing"
	pacingv1 "github.com/tatayooo/review-poc-go/internal/api/grpc/pacing/v1"
	psvc "github.com/tatayooo/review-poc-go/internal/domain/pacing/service"
	"github.com/tatayooo/review-poc-go/pkg/logger"
)

// PacingHandler adapts pacing decisions to gRPC.
type PacingHandler struct {
	pacingv1.UnimplementedPacingServiceServer
	svc  *psvc.PacingService
}

// NewPacing constructs the handler.
func NewPacing(svc *psvc.PacingService) *PacingHandler { return &PacingHandler{svc: svc} }

// Decide maps a decision request through the domain.
func (h *PacingHandler) Decide(ctx context.Context, req *pacingv1.DecideRequest) (*pacingv1.DecideResponse, error) {
	amount, err := h.svc.Decide(ctx, req.GetCampaignId())
	if err != nil {
		logger.Error(ctx, "pacing decide failed", "campaign", req.GetCampaignId(), "err", err)
		return nil, mapPacingError(err)
	}
	return &pacingv1.DecideResponse{AllowMinor: amount}, nil
}

// RegisterPlan validates and installs dayparts.
func (h *PacingHandler) RegisterPlan(ctx context.Context, req *pacingv1.RegisterPlanRequest) (*pacingv1.RegisterPlanResponse, error) {
	ws := make([]pacing.Window, 0, len(req.GetWindows()))
	for _, w := range req.GetWindows() {
		ws = append(ws, pacing.Window{
			Start: time.Duration(w.GetStartHour()) * time.Hour,
			End:   time.Duration(w.GetEndHour()) * time.Hour,
		})
	}
	if err := h.svc.RegisterPlan(ctx, req.GetCampaignId(), req.GetDailyCapMinor(), ws); err != nil {
		return nil, mapPacingError(err)
	}
	return &pacingv1.RegisterPlanResponse{Ok: true}, nil
}

// mapPacingError converts domain errors to transport codes (stubbed).
func mapPacingError(err error) error { return err }

var _ = campaignv1.Campaign{} // keep import shape honest
