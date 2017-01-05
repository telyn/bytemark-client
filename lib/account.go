package lib

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
)

// Account represents both the BigV and bmbilling parts of an account.
type Account struct {
	Name             string          `json:"name"`
	Owner            *billing.Person `json:"owner"`
	TechnicalContact *billing.Person `json:"technical_contact"`
	BillingID        int             `json:"billing_id"`
	BrainID          int             `json:"brain_id"`
	CardReference    string          `json:"card_reference"`
	Groups           []*brain.Group  `json:"groups"`
	Suspended        bool            `json:"suspended"`

	IsDefaultAccount bool `json:"-"`
}

func (a *Account) fillBrain(b *brain.Account) {
	if b != nil {
		a.BrainID = b.ID
		a.Groups = b.Groups
		a.Suspended = b.Suspended
		a.Name = b.Name
	}
}
func (a *Account) fillBilling(b *billing.Account) {
	if b != nil {
		a.BillingID = b.ID
		a.Owner = b.Owner
		a.TechnicalContact = b.TechnicalContact
		a.CardReference = b.CardReference
		a.Name = b.Name
	}
}

// CountVirtualMachines returns the number of virtual machines across all the Account's Groups.
func (a Account) CountVirtualMachines() (servers int) {
	for _, g := range a.Groups {
		servers += len(g.VirtualMachines)
	}
	return
}

// billingAccount copies all the billing parts of the account into a new billingAccount.
func (a Account) billingAccount() (b *billing.Account) {
	b = new(billing.Account)
	b.ID = a.BillingID
	b.Owner = a.Owner
	b.TechnicalContact = a.TechnicalContact
	b.CardReference = a.CardReference
	b.Name = a.Name
	return
}

func (pp Account) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return prettyprint.Run(wr, accountsTemplate, "account"+string(detail), pp)
}
