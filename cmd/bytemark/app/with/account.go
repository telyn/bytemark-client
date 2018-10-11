package with

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// Account gets an account name from a flag, then the account details from the API, then stitches it to the context
func Account(flagName string) func(*app.Context) error {
	return func(c *app.Context) (err error) {
		err = preflight(c)
		if err != nil {
			return
		}
		accName := c.String(flagName)
		c.Debug("flagName: %s accName: %s\n", flagName, accName)
		if accName == "" {
			accName = c.Config().GetIgnoreErr("account")
		}
		c.Debug("flagName: %s a4tName: %s\n", flagName, accName)

		acc, err := c.Client().GetAccount(accName)
		c.Account = &acc

		if err != nil {
			if _, ok := err.(lib.BillingAccountNotFound); !ok {
				return
			}
			err = nil
		}
		return
	}
}
