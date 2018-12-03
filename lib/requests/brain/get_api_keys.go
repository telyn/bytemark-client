package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// GetAPIKeys gets all API keys that you can currently see.
// In general this means those for your user, but users with cluster_admin
// will be able to see all API keys on the cluster
func GetAPIKeys(client lib.Client) (apiKeys []brain.APIKey, err error) {
	r, err := client.BuildRequest("GET", lib.BrainEndpoint, "/api_keys?view=overview")
	if err != nil {
		return
	}
	_, _, err = r.MarshalAndRun(nil, &apiKeys)
	return
}
