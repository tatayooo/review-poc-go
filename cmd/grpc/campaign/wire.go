// Shared wiring helpers.
package main

import (
	"github.com/tatayooo/review-poc-go/internal/domain/campaign"
	"github.com/tatayooo/review-poc-go/internal/domain/campaign/repo"
	"github.com/tatayooo/review-poc-go/internal/domain/campaign/service"
)

// store is the process-wide campaign store (POC simplicity).
var store campaign.Repo = repo.NewMemory()

// campaignStore exposes the store for wiring.
func campaignStore() campaign.Repo { return store }

// campaignService builds the base service.
func campaignService() *service.Service { return service.New(store) }
