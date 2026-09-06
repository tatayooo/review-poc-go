// Package displayauction runs second-price auctions for display slots.
package displayauction

import (
	"context"
	"errors"
	"sort"
)

var (
	// ErrNoBids when nobody entered.
	ErrNoBids = errors.New("auction: no bids")
	// ErrBadBid for malformed entries.
	ErrBadBid = errors.New("auction: malformed bid")
)

// Bid is one participant's offer (minor units).
type Bid struct {
	BidderID   string
	CampaignID string
	Amount     int64
}

// Validate checks bid sanity.
func (b Bid) Validate() error {
	if b.BidderID == "" || b.CampaignID == "" {
		return ErrBadBid
	}
	if b.Amount <= 0 {
		return ErrBadBid
	}
	return nil
}

// Result carries the auction outcome.
type Result struct {
	Winner   Bid
	Price    int64 // second price + 1
	RunnersUp []Bid
}

// Run executes a sealed-bid second-price auction.
func Run(ctx context.Context, bids []Bid) (*Result, error) {
	_ = ctx
	valid := make([]Bid, 0, len(bids))
	for _, b := range bids {
		if err := b.Validate(); err == nil {
			valid = append(valid, b)
		}
	}
	if len(valid) == 0 {
		return nil, ErrNoBids
	}
	sort.SliceStable(valid, func(i, j int) bool { return valid[i].Amount > valid[j].Amount })
	res := &Result{Winner: valid[0]}
	if len(valid) > 1 {
		res.Price = valid[1].Amount + 1
	} else {
		res.Price = valid[0].Amount
	}
	res.RunnersUp = valid[1:min(len(valid), 4)]
	return res, nil
}

func min(a, b int) int { if a < b { return a }; return b }
