package lib

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
