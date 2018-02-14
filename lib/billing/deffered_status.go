package billing

// DefferedStatus represents the ID and Deffered status of an account that is returned from bmbilling and
// at the moment, we are only intrested in the ID of this, as it converts a username into a billingID.
type DefferedStatus struct {
	ID       int  `json:"id,omitempty"`
	Deffered bool `json:"deffered,omitempty"`
}
