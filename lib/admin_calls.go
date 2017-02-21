package lib

import (
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func (c *bytemarkClient) GetVLANs() (vlans []*brain.VLAN, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/vlans")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &vlans)
	return
}

func (c *bytemarkClient) GetIPRange(id int) (ipRange *brain.IPRange, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/ip_ranges/%s", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &ipRange)
	return
}
