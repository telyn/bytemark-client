package brain

import (
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// RevokeAPIKey takes an API key id and revokes it.
func RevokeAPIKey(client lib.Client, id int) (err error) {
	r, err := client.BuildRequest("PUT", lib.BrainEndpoint, "/api_keys/%s", strconv.Itoa(id))
	if err != nil {
		return
	}

	apiKey := brain.APIKey{
		ExpiresAt: "00:00:00",
	}
	_, _, err = r.MarshalAndRun(apiKey, nil)
	return
}
