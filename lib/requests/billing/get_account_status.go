package billing

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
)

// GetAccountStatus gets the ID and Deffered status of a specified account
// at the moment, we are only intrested in the ID of this, as it converts a username into a billingID.
func GetAccountStatus(client lib.Client, username string) (account billing.DefferedStatus, err error) {
	req, err := client.BuildRequest("GET", lib.BillingEndpoint, "/api/v1/accounts/status?username=%s", username)
	if err != nil {
		return
	}
	accounts := []billing.DefferedStatus{}

	_, _, err = req.Run(nil, &accounts)

	if len(accounts) == 0 {
		err = fmt.Errorf("No accounts were returned with the username %s", username)
		return
	}
	account = accounts[0]
	return
}
