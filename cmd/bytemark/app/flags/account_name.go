package flags

import "github.com/BytemarkHosting/bytemark-client/lib"

// AccountName is used for all --account flags, excluding the global one.
type AccountName struct {
	AccountName        string
	Value              string
	SetFromCommandLine bool
}

// Set runs lib.ParseAccountName to make sure we get just the 'pure' account name; no cluster / endpoint details
func (name *AccountName) Set(value string) error {
	name.Value = value
	name.SetFromCommandLine = true
	return nil
}

// Preprocess sets the value of this flag to the global account flag if it's unset,
// and then runs lib.ParseAccountName
func (name *AccountName) Preprocess(c *Context) (err error) {
	if name.Value == "" {
		name.Value = c.Context.GlobalString("account")
	}
	name.AccountName = lib.ParseAccountName(name.Value, c.Config().GetIgnoreErr("account"))
	return
}

// String returns the AccountName as a string.
func (name AccountName) String() string {
	if name.AccountName == "" {
		return name.Value
	}
	return name.AccountName
}
