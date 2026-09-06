// Package handler adapts the campaign service to gRPC-shaped calls.
package handler

import (
	"context"

	campaignv1 "github.com/tatayooo/review-poc-go/internal/api/grpc/campaign/v1"
	"github.com/tatayooo/review-poc-go/internal/domain/campaign"
	"github.com/tatayooo/review-poc-go/internal/domain/campaign/service"
	"github.com/tatayooo/review-poc-go/pkg/logger"
)

// Handler exposes campaign operations over gRPC.
type Handler struct {
	campaignv1.UnimplementedCampaignServiceServer
	svc *service.Service
}

// New constructs a Handler.
func New(svc *service.Service) *Handler { return &Handler{svc: svc} }

// GetCampaign maps domain results to transport shapes.
func (h *Handler) GetCampaign(ctx context.Context, req *campaignv1.GetCampaignRequest) (*campaignv1.Campaign, error) {
	c, err := h.svc.Get(ctx, req.GetId())
	if err != nil {
		logger.Error(ctx, "get campaign failed", "err", err)
		return nil, mapError(err)
	}
	return toProto(c), nil
}

// toProto copies the aggregate into the transport shape.
func toProto(c *campaign.Campaign) *campaignv1.Campaign {
	return &campaignv1.Campaign{
		Id:           c.ID,
		Name:         c.Name,
		AdvertiserId: c.AdvertiserID,
		Status:       string(c.Status),
		Budget:       c.Budget,
	}
}
