package flags

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// AccountNameFlag is used for all --account flags, excluding the global one.
type AccountNameFlag struct {
	AccountName        string
	Value              string
	SetFromCommandLine bool
}

// Set runs lib.ParseAccountName to make sure we get just the 'pure' account name; no cluster / endpoint details
func (name *AccountNameFlag) Set(value string) error {
	name.Value = value
	name.SetFromCommandLine = true
	return nil
}

// Preprocess sets the value of this flag to the global account flag if it's unset,
// and then runs lib.ParseAccountName
func (name *AccountNameFlag) Preprocess(c *app.Context) (err error) {
	if name.Value == "" {
		name.Value = c.Context.GlobalString("account")
	}
	name.AccountName = lib.ParseAccountName(name.Value, c.Config().GetIgnoreErr("account"))
	return
}

// String returns the AccountNameFlag as a string.
func (name AccountNameFlag) String() string {
	if name.AccountName == "" {
		return name.Value
	}
	return name.AccountName
}
