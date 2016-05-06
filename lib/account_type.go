package lib

// brainAccount represents an account that's used by the brains.
type brainAccount struct {
	Name string `json:"name"`

	// the following cannot be set
	ID        int      `json:"id"`
	Suspended bool     `json:"suspended"`
	Groups    []*Group `json:"groups"`
}

// billingAccount represents the parts of an account that are discussed with bmbilling
type billingAccount struct {
	ID                 int     `json:"id"`
	Name               string  `json:"bigv_account_subscription"`
	Owner              *Person `json:"owner"`
	TechnicalContact   *Person `json:"tech"`
	OwnerID            int     `json:"owner_id" omitempty`
	CardReference      string  `json:"card_reference" omitempty`
	TechnicalContactID int     `json:"tech_id" omitempty`
}

type Account struct {
	Name             string   `json:"name"`
	Owner            *Person  `json:"owner"`
	TechnicalContact *Person  `json:"technical_contact"`
	BillingID        int      `json:"billing_id"`
	BrainID          int      `json:"brain_id"`
	CardReference    string   `json:"card_reference"`
	Groups           []*Group `json:"groups"`
	Suspended        bool     `json:"suspended"`
}

func (a *Account) FillBrain(b *brainAccount) {
	if b != nil {
		a.BrainID = b.ID
		a.Groups = b.Groups
		a.Suspended = b.Suspended
		a.Name = b.Name
	}
}
func (a *Account) FillBilling(b *billingAccount) {
	if b != nil {
		a.BillingID = b.ID
		a.Owner = b.Owner
		a.TechnicalContact = b.TechnicalContact
		a.CardReference = b.CardReference
		a.Name = b.Name
	}
}

func (a *Account) CountVirtualMachines() (servers int) {
	for _, g := range a.Groups {
		servers += len(g.VirtualMachines)
	}
	return
}

/*
func (a *Account) ToBillingAccount() *billingAccount {

}
func (a *Account) ToBrainAccount() *brainAccount {

}
*/
