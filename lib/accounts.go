package lib

import (
	"fmt"
)

// GetAccount takes an account name and returns a filled-out Account object
func (bigv *Client) GetAccount(name string) (account *Account, err error) {
	account = new(Account)
	path := fmt.Sprintf("/accounts/%s", name)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
