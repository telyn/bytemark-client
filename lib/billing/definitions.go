package billing

// Definitions is admin-modifiable parameters of bmbilling
type Definitions struct {
	TrialDays  int `json:"trial_days,omitempty"`
	TrialPence int `json:"trial_pence,omitempty"`
}
