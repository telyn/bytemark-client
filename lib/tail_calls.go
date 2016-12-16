package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func (c *bytemarkClient) GetTails() (tails brain.Tails, err error) {
	req, err := c.BuildRequest("GET", BrainEndpoint, "/admin/tails")
	if err != nil {
		return
	}

	_, _, err = req.Run(nil, &tails)
	return
}
