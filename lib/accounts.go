package lib

import (
	"bytes"
	"encoding/json"
)

/*func (c *Client) RegisterAccount() {

}*/

// getBillingAccount gets the billing account with the given name.
// Due to the way the billing API is implemented this is done by grabbing them all and looping *shrug*
func (c *bytemarkClient) getBillingAccount(name string) (account *billingAccount, err error) {
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
func (c *bytemarkClient) getBillingAccounts() (accounts []*billingAccount, err error) {
	req, err := c.BuildRequest("GET", EP_BILLING, "/api/v1/accounts")
	if err != nil {
		return
	}
	_, _, err = req.Run(nil, &accounts)
	return
}

// getBrainAccount gets the brain account with the given name.
func (c *bytemarkClient) getBrainAccount(name string) (account *brainAccount, err error) {
	err = c.validateAccountName(&name)
	if err != nil {
		return
	}
	account = new(brainAccount)

	req, err := c.BuildRequest("GET", EP_BRAIN, "/accounts/%s?view=overview&include_deleted=true", name)
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
	req, err := c.BuildRequestNoAuth("POST", EP_BILLING, "/api/v1/accounts")
	if err != nil {
		return nil, err
	}

	js, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	_, _, err = req.Run(bytes.NewBuffer(js), newAcc)
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
	account.FillBrain(brainAccount)
	account.FillBilling(billingAccount)

	return

}

func (c *bytemarkClient) getBrainAccounts() (accounts []*brainAccount, err error) {
	accounts = make([]*brainAccount, 1, 1)

	req, err := c.BuildRequest("GET", EP_BRAIN, "/accounts")
	if err != nil {
		return
	}

	_, _, err = req.Run(nil, &accounts)
	if err != nil {
		return
	}

	return
}

// Gets all Accounts you can see, merging data from both the brain and the billing
func (c *bytemarkClient) GetAccounts() (accounts []*Account, err error) {
	by_name := make(map[string]*Account)
	billingAccounts, err := c.getBillingAccounts()
	if err != nil {
		return
	}
	brainAccounts, err := c.getBrainAccounts()
	if err != nil {
		return
	}

	for _, b := range brainAccounts {
		if by_name[b.Name] == nil {
			by_name[b.Name] = new(Account)
		}
		by_name[b.Name].FillBrain(b)
	}
	for _, b := range billingAccounts {
		if by_name[b.Name] == nil {
			by_name[b.Name] = new(Account)
		}
		by_name[b.Name].FillBilling(b)
	}
	for _, a := range by_name {
		accounts = append(accounts, a)
	}
	return

}
