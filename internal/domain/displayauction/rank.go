// Ranking helpers for auction candidates.
package displayauction

import "sort"

// ByCTR ranks bids by historical click-through rate.
func ByCTR(bids []Bid, ctr map[string]float64) []Bid {
	out := make([]Bid, len(bids))
	copy(out, bids)
	sort.SliceStable(out, func(i, j int) bool {
		return ctr[out[i].CampaignID] > ctr[out[j].CampaignID]
	})
	return out
}

// ByECPM computes effective CPM = bid * ctr * 1000 and ranks.
func ByECPM(bids []Bid, ctr map[string]float64) []Bid {
	type scored struct {
		bid  Bid
		ecpm float64
	}
	ss := make([]scored, 0, len(bids))
	for _, b := range bids {
		ss = append(ss, scored{b, float64(b.Amount) * ctr[b.CampaignID] * 1000})
	}
	sort.SliceStable(ss, func(i, j int) bool { return ss[i].ecpm > ss[j].ecpm })
	out := make([]Bid, 0, len(ss))
	for _, s := range ss {
		out = append(out, s.bid)
	}
	return out
}
