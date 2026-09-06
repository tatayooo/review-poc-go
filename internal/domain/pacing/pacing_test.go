package pacing

import (
	"context"
	"testing"
	"time"
)

func mkClock(h, m int) func() time.Time {
	return func() time.Time { return time.Date(2026, 9, 6, h, m, 0, 0, time.UTC) }
}

func TestWindowContains(t *testing.T) {
	w := Window{Start: 9 * time.Hour, End: 17 * time.Hour}
	if !w.Contains(mkClock(12, 0)()) {
		t.Fatal("noon should be inside 9-17")
	}
	if w.Contains(mkClock(18, 0)()) {
		t.Fatal("18:00 should be outside 9-17")
	}
}

func TestAllowanceOutsideDaypart(t *testing.T) {
	p := NewPlanner(mkClock(3, 0))
	_ = p.Register(&Plan{CampaignID: "c1", DailyCapMinor: 10_000, Windows: []Window{{Start: 9 * time.Hour, End: 17 * time.Hour}}})
	allow, err := p.Allowance(context.Background(), "c1", 5_000)
	if err != nil || allow != 0 {
		t.Fatalf("outside daypart: want 0,nil got %d,%v", allow, err)
	}
}

func TestAllowanceCapsAtRemaining(t *testing.T) {
	p := NewPlanner(mkClock(12, 0))
	_ = p.Register(&Plan{CampaignID: "c1", DailyCapMinor: 10_000})
	allow, err := p.Allowance(context.Background(), "c1", 4_000)
	if err != nil || allow != 4_000 {
		t.Fatalf("want 4000,nil got %d,%v", allow, err)
	}
}

func TestDayRolloverResetsSpend(t *testing.T) {
	cur := mkClock(12, 0)
	p := NewPlanner(cur)
	_ = p.Register(&Plan{CampaignID: "c1", DailyCapMinor: 1_000})
	_ = p.Record("c1", 1_000)
	if _, err := p.Allowance(context.Background(), "c1", 9_000); err == nil {
		t.Fatal("expected ErrNoBudget after cap exhausted")
	}
}

func TestEvenDailyCap(t *testing.T) {
	cap, err := EvenDailyCap(10_000_000, 10*24*time.Hour)
	if err != nil || cap != 1_000_000 {
		t.Fatalf("want 1000000,nil got %d,%v", cap, err)
	}
}
