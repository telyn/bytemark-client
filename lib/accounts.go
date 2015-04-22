package lib

import (
	"fmt"
)

// GetAccount takes an account name and returns a filled-out Account object
func (bigv *BigVClient) GetAccount(name string) (account *Account, err error) {
	account = new(Account)
	path := fmt.Sprintf("/accounts/%s", name)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", account)
	if err != nil {
		switch err.(type) {
		case NotFoundError:
			newErr := err.(NotFoundError)
			newErr.Thing = name
			newErr.ThingType = "account"
		case NotAuthorizedError:
			newErr := err.(NotAuthorizedError)
			newErr.Thing = name
			newErr.ThingType = "account"
		}
	}

	return account, nil
}
