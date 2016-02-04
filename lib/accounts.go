package lib

// GetAccount takes an account name or ID and returns a filled-out Account object
func (c *bytemarkClient) GetAccount(name string) (account *Account, err error) {
	err = c.validateAccountName(&name)
	if err != nil {
		return
	}
	account = new(Account)

	req, err := c.BuildRequest("GET", EP_BRAIN, "/accounts/%s?view=overview", name)
	if err != nil {
		return
	}

	_, _, err = req.Run(nil, account)

	return
}

func (c *bytemarkClient) GetAccounts() (accounts []*Account, err error) {
	accounts = make([]*Account, 1, 1)

	req, err := c.BuildRequest("GET", EP_BRAIN, "/accounts")
	if err != nil {
		return
	}

	_, _, err = req.Run(nil, accounts)
	if err != nil {
		return
	}

	return

}
