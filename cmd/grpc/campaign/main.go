// Command campaign runs the campaign gRPC server (POC stub).
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tatayooo/review-poc-go/internal/api/grpc/handler"
	"github.com/tatayooo/review-poc-go/internal/domain/campaign/repo"
	"github.com/tatayooo/review-poc-go/internal/domain/campaign/service"
	"github.com/tatayooo/review-poc-go/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := repo.NewMemory()
	svc := service.New(store)
	h := handler.New(svc)
	logger.Info(ctx, "campaign grpc starting", "handlers", 1)
	_ = h // server wiring elided for the POC
}
