package lib

import "fmt"

// GetAccount takes an account name or ID and returns a filled-out Account object
func (bigv *bigvClient) GetAccount(name string) (*Account, error) {
	fmt.Println(name)
	err := bigv.validateAccountName(&name)
	fmt.Println(name)
	if err != nil {
		return nil, err
	}
	account := new(Account)

	path := BuildURL("/accounts/%s", name)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
