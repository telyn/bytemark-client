package lib

import (
	"encoding/json"
	"fmt"
)

func (bigv *Client) GetAccount(name string) (account *Account, err error) {
	account = new(Account)
	path := fmt.Sprintf("/accounts/%s", name)
	data, err := bigv.Request("GET", path, "")

	if err != nil {
		//TODO(telyn): good error handling here
		panic("Couldn't make request")
	}

	err = json.Unmarshal(data, account)
	if err != nil {
		fmt.Printf("Data returned was not an Account\r\n")
		fmt.Printf("%+v\r\n", account)

		return nil, err
	}
	return account, nil
}
