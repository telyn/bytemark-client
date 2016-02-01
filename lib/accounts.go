package lib

// GetAccount takes an account name or ID and returns a filled-out Account object
func (c *bytemarkClient) GetAccount(name string) (*Account, error) {
	err := c.validateAccountName(&name)
	if err != nil {
		return nil, err
	}
	account := new(Account)

	path := BuildURL("/accounts/%s?view=overview", name)

	err = c.RequestAndUnmarshal(true, "GET", path, "", account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (c *bytemarkClient) GetAccounts() ([]*Account, error) {
	accounts := make([]*Account, 1, 1)

	path := BuildURL("/accounts")

	err := c.RequestAndUnmarshal(true, "GET", path, "", &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil

}
