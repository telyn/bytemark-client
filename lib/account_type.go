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
	ID                 int     `json:"id,omitempty"`
	Name               string  `json:"bigv_account_subscription,omitempty"`
	Owner              *Person `json:"owner,omitempty"`
	TechnicalContact   *Person `json:"tech,omitempty"`
	OwnerID            int     `json:"owner_id,omitempty"`
	CardReference      string  `json:"card_reference,omitempty"`
	TechnicalContactID int     `json:"tech_id,omitempty"`
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

func (a *Account) fillBrain(b *brainAccount) {
	if b != nil {
		a.BrainID = b.ID
		a.Groups = b.Groups
		a.Suspended = b.Suspended
		a.Name = b.Name
	}
}
func (a *Account) fillBilling(b *billingAccount) {
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

func (a *Account) billingAccount() (b *billingAccount) {
	b = new(billingAccount)
	b.ID = a.BillingID
	b.Owner = a.Owner
	b.TechnicalContact = a.TechnicalContact
	b.CardReference = a.CardReference
	b.Name = a.Name
	return
}

/*
func (a *Account) ToBillingAccount() *billingAccount {

}
func (a *Account) ToBrainAccount() *brainAccount {

}
*/

type Person struct {
	ID          int    `json:"id,omitempty"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	BackupEmail string `json:"email_backup,omitempty"`

	// only set in the creation request
	Password string `json:"password"`

	FirstName   string `json:"firstname"`
	LastName    string `json:"surname"`
	Address     string `json:"address"`
	City        string `json:"city"`
	StateCounty string `json:"statecounty,omitempty"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	MobilePhone string `json:"phonemobile,omitempty"`

	Organization         string `json:"organization,omitempty"`
	OrganizationDivision string `json:"division,omitempty"`
	VATNumber            string `json:"vatnumber,omitempty"`
}
