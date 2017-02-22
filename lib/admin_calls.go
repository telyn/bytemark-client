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

func (c *bytemarkClient) GetIPRanges() (ipRanges []*brain.IPRange, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/ip_ranges")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &ipRanges)
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

func (c *bytemarkClient) GetHeads() (heads []*brain.Head, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/heads")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &heads)
	return
}

func (c *bytemarkClient) GetHead(id int) (head *brain.Head, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/heads/%s", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &head)
	return
}

func (c *bytemarkClient) GetStoragePools() (storagePools []*brain.StoragePool, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/storage_pools")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &storagePools)
	return
}

func (c *bytemarkClient) GetStoragePool(idOrLabel string) (storagePool *brain.StoragePool, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/storage_pools/%s", idOrLabel)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &storagePool)
	return
}

func (c *bytemarkClient) GetMigratingVMs() (vms []*brain.VirtualMachine, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/migrating_vms")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &vms)
	return
}
