package billing

// Account represents the parts of an account that are discussed with bmbilling
type Account struct {
	ID                 int     `json:"id,omitempty"`
	Name               string  `json:"bigv_account_subscription,omitempty"`
	Owner              *Person `json:"owner,omitempty"`
	TechnicalContact   *Person `json:"tech,omitempty"`
	OwnerID            int     `json:"owner_id,omitempty"`
	CardReference      string  `json:"card_reference,omitempty"`
	TechnicalContactID int     `json:"tech_id,omitempty"`
}
