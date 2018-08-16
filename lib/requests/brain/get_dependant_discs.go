package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// Function to get array of discs on specified tail
func GetDiscsOnTail(client lib.Client, id string, at string) (servers brain.Discs, err error) {
	var r lib.Request

	if at != "" {
		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/tails/%s/discs?at=%s", id, at)
	} else {
		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/tails/%s/discs", id)
	}

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &servers)

	return
}

// Function to get array of discs on specified storage pool
func GetDiscsOnStoragePool(client lib.Client, id string, at string) (servers brain.Discs, err error) {
	var r lib.Request

	if at != "" {
		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/storage_pools/%s/discs?at=%s", id, at)
	} else {
		r, err = client.BuildRequest("GET", lib.BrainEndpoint, "/admin/storage_pools/%s/discs", id)
	}

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &servers)

	return
}
