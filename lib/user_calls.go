package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// Getbigv.User grabs the named user from the brain
func (c *bytemarkClient) GetUser(name string) (user brain.User, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/users/%s", name)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &user)
	if err != nil {
		return
	}
	return
}
