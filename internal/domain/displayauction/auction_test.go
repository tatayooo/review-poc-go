package displayauction

import (
	"context"
	"testing"
)

func TestSecondPrice(t *testing.T) {
	bids := []Bid{
		{"b1", "c1", 500},
		{"b2", "c2", 700},
		{"b3", "c3", 300},
	}
	r, err := Run(context.Background(), bids)
	if err != nil {
		t.Fatal(err)
	}
	if r.Winner.BidderID != "b2" {
		t.Fatalf("want b2 got %s", r.Winner.BidderID)
	}
	if r.Price != 501 {
		t.Fatalf("second price want 501 got %d", r.Price)
	}
}

func TestSingleBidPaysOwn(t *testing.T) {
	r, err := Run(context.Background(), []Bid{{"b1", "c1", 900}})
	if err != nil || r.Price != 900 {
		t.Fatalf("want 900,nil got %d,%v", r.Price, err)
	}
}

func TestInvalidBidsSkipped(t *testing.T) {
	r, err := Run(context.Background(), []Bid{
		{"", "c1", 100},
		{"b2", "c2", 0},
		{"b3", "c3", 200},
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.Winner.BidderID != "b3" {
		t.Fatalf("want b3 got %s", r.Winner.BinnerID())
	}
}

func TestByECPM(t *testing.T) {
	bids := []Bid{{"b1", "c1", 1000}, {"b2", "c2", 500}}
	ctr := map[string]float64{"c1": 0.001, "c2": 0.01}
	got := ByECPM(bids, ctr)
	if got[0].CampaignID != "c2" {
		t.Fatalf("c2 should outrank on eCPM, got %s", got[0].CampaignID)
	}
}
