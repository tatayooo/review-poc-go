// Package campaignv1 holds minimal generated-shaped types for the POC.
package campaignv1

// Campaign is the transport shape for a campaign.
type Campaign struct {
	Id           string
	Name         string
	AdvertiserId string
	Status       string
	Budget       int64
}

// GetCampaignRequest selects a campaign by id.
type GetCampaignRequest struct{ Id string }

// GetId returns the request id.
func (r *GetCampaignRequest) GetId() string { return r.Id }

// UnimplementedCampaignServiceServer keeps forward compat with real protos.
type UnimplementedCampaignServiceServer struct{}

// mapError placeholder lives in handler to avoid import cycles.
