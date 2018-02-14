package billing

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
)

// GetAccountDeferredStatus gets the ID and Deferred status of a specified account
// at the moment, we are only intrested in the ID of this, as it converts a username into a billingID.
func GetAccountDeferredStatus(client lib.Client, username string) (account billing.DeferredStatus, err error) {
	req, err := client.BuildRequest("GET", lib.BillingEndpoint, "/api/v1/accounts/status?username=%s", username)
	if err != nil {
		return
	}
	accounts := []billing.DeferredStatus{}

	_, _, err = req.Run(nil, &accounts)

	if len(accounts) == 0 {
		err = fmt.Errorf("No accounts were returned with the username %s", username)
		return
	}
	account = accounts[0]
	return
}
