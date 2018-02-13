package billing

// Account represents the parts of an account that are discussed with bmbilling
type DefferedStatus struct {
	ID       int  `json:"id,omitempty"`
	Deffered bool `json:"deffered,omitempty"`
}
