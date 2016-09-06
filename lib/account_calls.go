package lib

import (
	"bytes"
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

/*func (c *Client) RegisterAccount() {

}*/

// getBillingAccount gets the billing account with the given name.
// Due to the way the billing API is implemented this is done by grabbing them all and looping *shrug*
func (c *bytemarkClient) getBillingAccount(name string) (account *billing.Account, err error) {
	accounts, err := c.getBillingAccounts()
	if err != nil {
		return
	}
	for _, account = range accounts {
		if account.Name == name {
			return
		}
	}
	account = nil
	return
}

// getBillingAccounts returns all the billing accounts the currently logged in user can see.
func (c *bytemarkClient) getBillingAccounts() (accounts []*billing.Account, err error) {
	if c.billingEndpoint == "" {
		return make([]*billing.Account, 0), nil
	}
	req, err := c.BuildRequest("GET", BillingEndpoint, "/api/v1/accounts")
	if err != nil {
		return
	}
	_, _, err = req.Run(nil, &accounts)
	return
}

// getBrainAccount gets the brain account with the given name.
func (c *bytemarkClient) getBrainAccount(name string) (account *brain.Account, err error) {
	err = c.validateAccountName(&name)
	if err != nil {
		return
	}
	account = new(brain.Account)

	req, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s?view=overview&include_deleted=true", name)
	if err != nil {
		return nil, err
	}

	_, _, err = req.Run(nil, account)

	return
}

/*
// CreateAccount creates a new account with you as the owner.
func (c *bytemarkClient) CreateAccount(account *Account) (newAccount *Account, err error) {
    // TODO(telyn): implement
	return nil, nil
}*/

// RegisterNewAccount registers a new account with bmbilling. This will create a new user for the owner.
// If you would like an extra account attached to your regular user, use CreateAccount
func (c *bytemarkClient) RegisterNewAccount(acc *Account) (newAcc *Account, err error) {
	req, err := c.BuildRequestNoAuth("POST", BillingEndpoint, "/api/v1/accounts")
	if err != nil {
		return nil, err
	}

	js, err := json.Marshal(acc.billingAccount())
	if err != nil {
		return nil, err
	}

	status, _, err := req.Run(bytes.NewBuffer(js), newAcc)
	if err != nil {
		if _, ok := err.(*json.InvalidUnmarshalError); !ok {
			return newAcc, err
		}
	}
	if status == 202 {
		return newAcc, AccountCreationDeferredError{}
	}
	return newAcc, err
}

// GetAccount takes an account name or ID and returns a filled-out Account object
func (c *bytemarkClient) GetAccount(name string) (account *Account, err error) {
	billingAccount, err := c.getBillingAccount(name)
	if err != nil {
		return nil, err
	}
	brainAccount, err := c.getBrainAccount(name)
	if err != nil {
		return nil, err
	}
	account = new(Account)
	account.fillBrain(brainAccount)
	account.fillBilling(billingAccount)

	return

}

func (c *bytemarkClient) getBrainAccounts() (accounts []*brain.Account, err error) {
	accounts = make([]*brain.Account, 1, 1)

	req, err := c.BuildRequest("GET", BrainEndpoint, "/accounts")
	if err != nil {
		return
	}

	_, _, err = req.Run(nil, &accounts)
	if err != nil {
		return
	}

	return
}

func (c *bytemarkClient) getDefaultBillingAccount() (*billing.Account, error) {
	if c.brainEndpoint == "https://int.bigv.io" {
		return &billing.Account{Name: "bytemark"}, nil
	}
	billAccs, err := c.getBillingAccounts()
	if err != nil {
		return nil, err
	}

	return billAccs[0], nil
}

// GetDefaultAccount gets the account *most likely* to be your default account.
// This is the first account returned by the billing endpoint,
// with the brain's data for it attached. Fingers crossed.
// Returns the default billing account with NoDefaultAccountError if there's
// not a bigv_subscription_account on the billing account, and returns nil
func (c *bytemarkClient) GetDefaultAccount() (*Account, error) {
	acc := new(Account)
	// there is only one account worth mentioning on int.
	if c.brainEndpoint == "https://int.bigv.io" {
		brainAcc, err := c.getBrainAccount("bytemark")
		if err != nil {
			return nil, err
		}
		acc.fillBrain(brainAcc)
		return acc, nil
	}
	billAcc, err := c.getDefaultBillingAccount()
	if err != nil {
		return nil, NoDefaultAccountError{err}
	}

	acc.fillBilling(billAcc)

	if billAcc.Name != "" {
		brainAcc, err := c.getBrainAccount(billAcc.Name)
		if err != nil {
			return nil, err
		}

		acc.fillBrain(brainAcc)
		return acc, nil
	}

	return acc, NoDefaultAccountError{}
}

// Gets all Accounts you can see, merging data from both the brain and the billing
func (c *bytemarkClient) GetAccounts() (accounts []*Account, err error) {
	byName := make(map[string]*Account)
	billingAccounts, err := c.getBillingAccounts()
	if err != nil {
		return
	}
	brainAccounts, err := c.getBrainAccounts()
	if err != nil {
		return
	}

	for _, b := range brainAccounts {
		if byName[b.Name] == nil {
			byName[b.Name] = new(Account)
		}
		byName[b.Name].fillBrain(b)
	}
	for _, b := range billingAccounts {
		if byName[b.Name] == nil {
			byName[b.Name] = new(Account)
		}
		byName[b.Name].fillBilling(b)
	}
	for _, a := range byName {
		accounts = append(accounts, a)
	}
	return

}

// Overview is a combination of a user's default account, their username, and all the accounts they have access to see.
type Overview struct {
	DefaultAccount *Account
	Username       string
	Accounts       []*Account
}

// GetOverview gets an Overview for everything the user can access at bytemark
func (c *bytemarkClient) GetOverview() (*Overview, error) {
	o := new(Overview)
	acc, err := c.GetDefaultAccount()
	if err != nil {
		return nil, err
	}
	accs, err := c.GetAccounts()
	if err != nil {
		return nil, err
	}
	o.DefaultAccount = acc
	o.Accounts = accs
	o.Username = c.authSession.Username
	return o, nil
}
