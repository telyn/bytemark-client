package billing

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
)

// getBillingAccountID gets the billing account ID from the given name.
func GetBillingAccountID(client lib.Client, username string) (accountID billing.DefferedStatus, err error) {
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
	accountID = accounts[0]
	return
}
