package billing

// SPPTokenRequest is the object sent to bmbilling to get a token to pass to SPP.
type SPPTokenRequest struct {
	// needs to be able to be nil, so pointer
	Owner      *Person `json:"owner,omitempty"`
	CardEnding string  `json:"card_ending"`
}
