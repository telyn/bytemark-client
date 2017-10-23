package billing

// Account represents the parts of an account that are discussed with bmbilling
type Account struct {
	ID                 int    `json:"id,omitempty"`
	Name               string `json:"bigv_account_subscription,omitempty"`
	Owner              Person `json:"owner,omitempty"`
	TechnicalContact   Person `json:"tech,omitempty"`
	OwnerID            int    `json:"owner_id,omitempty"`
	CardReference      string `json:"card_reference,omitempty"`
	TechnicalContactID int    `json:"tech_id,omitempty"`

	InvoiceTerms     int    `json:"invoice_terms,omitempty"`
	PaymentMethod    string `json:"payment_method,omitempty"`
	EarliestActivity string `json:"earliest_activity,omitempty"`
}

// IsValid returns true if the account is valid - whether its fields are set, or it should be considered a null account
func (a Account) IsValid() bool {
	return a.ID != 0 || a.Name != ""
}
