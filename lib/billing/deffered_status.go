package billing

// DeferredStatus represents the ID and Deferred status of an account that is returned from bmbilling and
// at the moment, we are only intrested in the ID of this, as it converts a username into a billingID.
type DeferredStatus struct {
	ID       int  `json:"id,omitempty"`
	Deferred bool `json:"deferred,omitempty"`
}
