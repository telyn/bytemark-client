package billing

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	types "github.com/BytemarkHosting/bytemark-client/lib/billing"
)

// UpdateDefinition sets the bmbilling definitions to defs (for all the ones that are non-zero, anyway)
func UpdateDefinition(client lib.Client, def types.Definition) (err error) {
	if def.Name == "" || def.Value == "" {
		return fmt.Errorf("Definition must have a non-blank name and value")
	}
	req, err := client.BuildRequest("PUT", lib.BillingEndpoint, "/api/v1/definitions/%s", def.Name)
	if err != nil {
		return
	}
	_, _, err = req.MarshalAndRun(def, nil)
	return
}
