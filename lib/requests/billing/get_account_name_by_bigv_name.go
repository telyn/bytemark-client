package billing

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
)

// GetAccountByBigVName gets a billing account by bigv account name
func GetAccountByBigVName(client lib.Client, bigvname string) (account billing.Account, err error) {
	req, err := client.BuildRequest("GET", lib.BillingEndpoint, "/api/v1/accounts?bigv_account_name=%s", bigvname)
	if err != nil {
		return
	}
	accounts := []billing.Account{}

	_, _, err = req.Run(nil, &accounts)

	if len(accounts) == 0 {
		err = fmt.Errorf("No accounts were returned with the BigV account name %s", bigvname)
		return
	}
	account = accounts[0]
	return
}
