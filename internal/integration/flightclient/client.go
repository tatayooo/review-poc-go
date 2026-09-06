// Package flightclient is the external flight-schedule adapter (stub).
package flightclient

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// ErrUnavailable wraps transport failures for retry mapping.
var ErrUnavailable = errors.New("flightclient: unavailable")

// Client fetches flight windows.
type Client struct {
	baseURL string
	http    *http.Client
}

// New builds a client with sane timeouts.
func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: 3 * time.Second}}
}

// Flight is one scheduled window from the external system.
type Flight struct {
	CampaignID string
	Start      time.Time
	End        time.Time
}

// Fetch pulls flights for a campaign; stubbed to fail-closed offline.
func (c *Client) Fetch(ctx context.Context, campaignID string) ([]Flight, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/flights?campaign="+campaignID, nil)
	if err != nil {
		return nil, ErrUnavailable
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, ErrUnavailable
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnavailable
	}
	return nil, ErrUnavailable // decoding elided for POC
}
