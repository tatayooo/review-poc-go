// Package pacingv1 holds generated-shaped pacing types for the POC.
package pacingv1

// DecideRequest asks for the current spendable amount.
type DecideRequest struct{ CampaignId string }

// GetCampaignId returns the id.
func (r *DecideRequest) GetCampaignId() string { return r.CampaignId }

// DecideResponse carries the allowance.
type DecideResponse struct{ AllowMinor int64 }

// RegisterPlanRequest installs a pacing plan.
type RegisterPlanRequest struct {
	CampaignId    string
	DailyCapMinor int64
	Windows       []*Window
}

// GetCampaignId returns the id.
func (r *RegisterPlanRequest) GetCampaignId() string { return r.CampaignId }

// GetDailyCapMinor returns the cap.
func (r *RegisterPlanRequest) GetDailyCapMinor() int64 { return r.DailyCapMinor }

// GetWindows returns the windows.
func (r *RegisterPlanRequest) GetWindows() []*Window { return r.Windows }

// Window is a daypart in hours.
type Window struct{ StartHour, EndHour int32 }

// GetStartHour returns the start.
func (w *Window) GetStartHour() int32 { return w.StartHour }

// GetEndHour returns the end.
func (w *Window) GetEndHour() int32 { return w.EndHour }

// RegisterPlanResponse acknowledges installation.
type RegisterPlanResponse struct{ Ok bool }

// UnimplementedPacingServiceServer keeps forward compat.
type UnimplementedPacingServiceServer struct{}
