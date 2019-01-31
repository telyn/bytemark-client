package lib

import "errors"

type AccountName string

// AccountPath returns an API path to the named account
func (a AccountName) AccountPath() (string, error) {
	if a == "" {
		return "", errors.New("Empty account name specified, cannot create a path. Make sure checkAccountPather was called first. This is a bug")
	}
	return "/accounts/" + string(a), nil
}
