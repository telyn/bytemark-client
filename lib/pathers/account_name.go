package pathers

import "errors"

// AccountName is the name for a Brain account - the same as is used in
// virtual-machine.group.account.uk0.bigv.io paths, for example.
type AccountName string

// AccountPath returns a BigV API path to the named account
func (a AccountName) AccountPath() (string, error) {
	if a == "" {
		return "", errors.New("Empty account name specified, cannot create a path. Make sure checkAccountPather was called first. This is a bug")
	}
	return "/accounts/" + string(a), nil
}
