package lib

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
)

/*func (c *Client) RegisterAccount() {

}*/

// getBillingAccount gets the billing account with the given name.
// Due to the way the billing API is implemented this is done by grabbing them all and looping *shrug*
func (c *bytemarkClient) getBillingAccount(name string) (account billing.Account, err error) {
	accounts, err := c.getBillingAccounts()
	if err != nil {
		return
	}
	for _, account = range accounts {
		if account.Name == name {
			return
		}
	}
	err = fmt.Errorf("Couldn't find a billing account called %s", name)
	return
}

// getBillingAccounts returns all the billing accounts the currently logged in user can see.
func (c *bytemarkClient) getBillingAccounts() (accounts []billing.Account, err error) {
	if c.urls.Billing == "" {
		return make([]billing.Account, 0), nil
	}
	req, err := c.BuildRequest("GET", BillingEndpoint, "/api/v1/accounts")
	if err != nil {
		return
	}
	oldfile := log.LogFile
	log.LogFile = nil
	_, _, err = req.Run(nil, &accounts)
	log.LogFile = oldfile
	return
}

// getBrainAccount gets the brain account with the given name.
func (c *bytemarkClient) getBrainAccount(name string) (account *brain.Account, err error) {
	err = c.EnsureAccountName(&name)
	if err != nil {
		return
	}

	req, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s?view=overview&include_deleted=true", name)
	if err != nil {
		return
	}

	_, _, err = req.Run(nil, &account)

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
func (c *bytemarkClient) RegisterNewAccount(acc Account) (newAcc Account, err error) {
	req, err := c.BuildRequestNoAuth("POST", BillingEndpoint, "/api/v1/accounts")
	if err != nil {
		return
	}

	// prevent password & card reference from being written to debug log
	// this is a bit of a sledgehammer
	// TODO make it not a sledgehammer somehow
	oldfile := log.LogFile
	log.LogFile = nil

	outputBillingAcc := billing.Account{}

	status, _, err := req.MarshalAndRun(acc.billingAccount(), &outputBillingAcc)
	if err != nil {
		if _, ok := err.(*json.InvalidUnmarshalError); !ok {
			return newAcc, err
		}
	}
	newAcc.fillBilling(outputBillingAcc)

	log.LogFile = oldfile
	if status == http.StatusAccepted {
		return newAcc, AccountCreationDeferredError{}
	}
	return newAcc, err
}

// GetAccount takes an account name or ID and returns a filled-out Account object
func (c *bytemarkClient) GetAccount(name string) (account Account, err error) {
	if name == "" {
		return c.GetDefaultAccount()
	}
	billingAccount, err := c.getBillingAccount(name)
	if err != nil {
		return Account{}, err
	}
	brainAccount, err := c.getBrainAccount(name)
	if err != nil {
		return Account{}, err
	}
	account.fillBrain(brainAccount)
	account.fillBilling(billingAccount)

	return

}

func (c *bytemarkClient) getBrainAccounts() (accounts []brain.Account, err error) {
	accounts = make([]brain.Account, 1)

	req, err := c.BuildRequest("GET", BrainEndpoint, "/accounts?view=overview&include_deleted=true")
	if err != nil {
		return
	}

	_, _, err = req.Run(nil, &accounts)
	if err != nil {
		return
	}

	return
}

// getDefaultBillingAccount gets the default billing account for this user.
func (c *bytemarkClient) getDefaultBillingAccount() (acc billing.Account, err error) {
	if c.urls.Billing == "https://int.bigv.io" {
		acc = billing.Account{Name: "bytemark"}
		return
	}
	billAccs, err := c.getBillingAccounts()
	if err != nil {
		return
	}
	if len(billAccs) == 0 {
		return
	}
	acc = billAccs[0]
	return
}

// GetDefaultAccount gets the account *most likely* to be your default account.
// This is the first account returned by the billing endpoint,
// with the brain's data for it attached. Fingers crossed.
// Returns the default billing account with NoDefaultAccountError if there's
// not a bigv_subscription_account on the billing account, and returns nil
func (c *bytemarkClient) GetDefaultAccount() (acc Account, err error) {
	acc.IsDefaultAccount = true
	billAcc, err := c.getDefaultBillingAccount()
	log.Debugf(log.LvlMisc, "billAcc: %#v, isValid: %v, err: %v\r\n", billAcc, billAcc.IsValid(), err)
	if err != nil {
		return
	}
	if billAcc.IsValid() {
		acc.fillBilling(billAcc)
		var brainAcc brain.Account

		brainAcc, err = c.getBrainAccount(billAcc.Name)
		if err != nil {
			return
		}
		acc.fillBrain(brainAcc)
	} else {
		var brainAccs []brain.Account

		brainAccs, err = c.getBrainAccounts()
		if err != nil {
			return
		}
		acc.fillBrain(brainAccs[0])
	}
	return
}

// Gets all Accounts you can see, merging data from both the brain and the billing
func (c *bytemarkClient) GetAccounts() (accounts Accounts, err error) {
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
		accounts = append(accounts, *a)
	}
	return

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
