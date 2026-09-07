package service

import (
	"context"
	"testing"

	"github.com/tatayooo/review-poc-go/internal/domain/creative"
)

type memRepo struct{ items map[string]*creative.Creative }

func (m *memRepo) Get(_ context.Context, id string) (*creative.Creative, error) {
	c, ok := m.items[id]
	if !ok {
		return nil, creative.ErrRejected
	}
	return c, nil
}
func (m *memRepo) ListByCampaign(_ context.Context, cid string) ([]*creative.Creative, error) {
	var out []*creative.Creative
	for _, c := range m.items {
		if c.CampaignID == cid {
			out = append(out, c)
		}
	}
	return out, nil
}
func (m *memRepo) Save(_ context.Context, c *creative.Creative) error { m.items[c.ID] = c; return nil }

func mkCreative(title string) *creative.Creative {
	return &creative.Creative{
		ID: "cr1", CampaignID: "c1", Title: title,
		LandingURL: "https://shop.example/landing",
		Assets:     []creative.Asset{{URL: "https://cdn.example/a.png", Format: creative.FormatImage, Aspect: creative.Aspect{W: 16, H: 9}, BytesMax: 1 << 20}},
	}
}

func TestSubmitBannedPhrase(t *testing.T) {
	s := NewModeration(&memRepo{items: map[string]*creative.Creative{}})
	if err := s.Submit(context.Background(), mkCreative("FREE MONEY now")); err == nil {
		t.Fatal("banned phrase must fail submit")
	}
}

func TestAutoModerateApproves(t *testing.T) {
	r := &memRepo{items: map[string]*creative.Creative{}}
	s := NewModeration(r)
	c := mkCreative("Summer Sale")
	if err := s.Submit(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	if err := s.AutoModerate(context.Background(), "cr1"); err != nil {
		t.Fatal(err)
	}
	got, _ := r.Get(context.Background(), "cr1")
	if got.State != creative.ReviewApproved {
		t.Fatalf("want APPROVED got %s", got.State)
	}
}

func TestAutoModerateRejectsVendorHTML(t *testing.T) {
	r := &memRepo{items: map[string]*creative.Creative{}}
	s := NewModeration(r)
	c := mkCreative("Fair deal")
	c.Assets = append(c.Assets, creative.Asset{URL: "https://doubleclick.example/x.html", Format: creative.FormatHTML, Aspect: creative.Aspect{W: 1, H: 1}, BytesMax: 4096})
	_ = s.Submit(context.Background(), c)
	if err := s.AutoModerate(context.Background(), "cr1"); err != nil {
		t.Fatal(err)
	}
	got, _ := r.Get(context.Background(), "cr1")
	if got.State != creative.ReviewRejected || got.RejectCode != "HTML_VENDOR" {
		t.Fatalf("want REJECTED/HTML_VENDOR got %s/%s", got.State, got.RejectCode)
	}
}
