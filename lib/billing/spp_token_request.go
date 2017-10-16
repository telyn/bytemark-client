package billing

type SPPTokenRequest struct {
	// needs to be able to be nil, so pointer
	Owner      *Person `json:"owner,omitempty"`
	CardEnding string  `json:"card_ending"`
}
