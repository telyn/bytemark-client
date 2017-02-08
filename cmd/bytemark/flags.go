package main

// AccountNameFlag is used for all --account flags, including the global one.
type AccountNameFlag string

// Set runs lib.Client.ParseAccountName using the global.Client to make sure we get just the 'pure' account name; no cluster / endpoint details
func (name *AccountNameFlag) Set(value string) error {
	*name = AccountNameFlag(global.Client.ParseAccountName(value, global.Config.GetIgnoreErr("account")))
	return nil
}

// String returns the AccountNameFlag as a string.
func (name *AccountNameFlag) String() string {
	return string(*name)
}
