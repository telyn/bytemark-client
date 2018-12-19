package brain

import (
	"fmt"
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib"
)

// DeleteAPIKey takes an API key id or label and revokes it.
func DeleteAPIKey(client lib.Client, id string) (err error) {
	if _, err := strconv.Atoi(id); err != nil {
		apikeys, err := GetAPIKeys(client)
		if err != nil {
			return err
		}
		found := false
		for _, k := range apikeys {
			if k.Label == id {
				id = strconv.Itoa(k.ID)
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Could not find an api key called %q", id)
		}
	}
	r, err := client.BuildRequest("DELETE", lib.BrainEndpoint, "/api_keys/%s", id)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}
