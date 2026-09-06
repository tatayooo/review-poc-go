// Pacing wiring for the campaign binary.
package main

import (
	"time"

	"github.com/tatayooo/review-poc-go/internal/api/grpc/handler"
	prepo "github.com/tatayooo/review-poc-go/internal/domain/pacing/repo"
	psvc "github.com/tatayooo/review-poc-go/internal/domain/pacing/service"
	"github.com/tatayooo/review-poc-go/internal/domain/pacing"
)

// wirePacing builds the pacing stack with real clock.
func wirePacing(h *handler.Handler) (*handler.PacingHandler, *prepo.Metrics) {
	planner := pacing.NewPlanner(time.Now)
	svc := psvc.NewPacing(campaignStore(), planner, time.Now)
	return handler.NewPacing(svc), prepo.NewMetrics()
}
