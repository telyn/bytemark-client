package flags

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// AccountNameFlag is used for all --account flags, excluding the global one.
type AccountNameFlag struct {
	// AccountName is the actual name that will be passed on to API calls, and
	// is made by checking the contents of Value are a valid account. If Value
	// is unset then the value of the 'account' config variable is used
	AccountName string
	// Value is the raw input to the flag, and can be used as the default when
	// creating the flag.
	Value string
	// SetFromCommandLine is false by default but is set to true when Set is
	// called. This allows setting a default by setting Value by yourself - Set
	// is called from urfave/cli's flag-parsing code.
	SetFromCommandLine bool
}

// Set sets Value and SetFromCommandLine on the flag
func (name *AccountNameFlag) Set(value string) error {
	name.Value = value
	name.SetFromCommandLine = true
	return nil
}

// Preprocess sets the value of this flag to the global account flag if it's
// unset, and then runs lib.ParseAccountName to set AccountName. This is an
// implementation of `app.Preprocessor`, which is detected and called
// automatically by actions created with `app.Action`
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
