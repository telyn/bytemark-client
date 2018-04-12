package billing

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	types "github.com/BytemarkHosting/bytemark-client/lib/billing"
)

// GetDefinition gets the named bmbilling definition
func GetDefinition(client lib.Client, name string) (def types.Definition, err error) {
	req, err := client.BuildRequest("GET", lib.BillingEndpoint, "/api/v1/definitions/%s", name)
	if err != nil {
		return
	}
	_, _, err = req.Run(nil, &def)
	return
}
