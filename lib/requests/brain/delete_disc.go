package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"

	"fmt"
)

// DeleteDiscByID takes a disc ID and removes the specified disc from the given virtual machine
func DeleteDiscByID(client lib.Client, discID string) (err error) {
	r, err := client.BuildRequest("DELETE", lib.BrainEndpoint, "/discs/%s?purge=true", discID)

	if err != nil {
		return
	}

	fmt.Println(r)

	_, _, err = r.Run(nil, nil)

	return
}