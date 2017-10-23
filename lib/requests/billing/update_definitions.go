package billing

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	types "github.com/BytemarkHosting/bytemark-client/lib/billing"
)

func UpdateDefinitions(client lib.Client, defs types.Definitions) (err error) {
	req, err := client.BuildRequest("PUT", lib.BillingEndpoint, "/api/v1/definitions")
	if err != nil {
		return
	}
	_, _, err = req.MarshalAndRun(defs, nil)
	return
}
