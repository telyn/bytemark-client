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

func (c *bytemarkClient) GetVLAN(num int) (vlan *brain.VLAN, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/vlans/%s", strconv.Itoa(num))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &vlan)
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

func (c *bytemarkClient) GetHead(idOrLabel string) (head *brain.Head, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/heads/%s", idOrLabel)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &head)
	return
}

func (c *bytemarkClient) GetTails() (tails []*brain.Tail, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/tails")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &tails)
	return
}

func (c *bytemarkClient) GetTail(idOrLabel string) (tail *brain.Tail, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/tails/%s", idOrLabel)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &tail)
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

func (c *bytemarkClient) GetStoppedEligibleVMs() (vms []*brain.VirtualMachine, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/stopped_eligible_vms")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &vms)
	return
}

func (c *bytemarkClient) GetRecentVMs() (vms []*brain.VirtualMachine, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/recent_vms")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &vms)
	return
}

func (c *bytemarkClient) MigrateDisc(disc int, newStoragePool string) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/discs/%s/migrate", strconv.Itoa(disc))
	if err != nil {
		return
	}

	params := map[string]string{}
	if newStoragePool != "" {
		params["new_pool_spec"] = newStoragePool
	}

	_, _, err = r.MarshalAndRun(params, nil)
	return
}

func (c *bytemarkClient) MigrateVirtualMachine(vmName *VirtualMachineName, newHead string) (err error) {
	vm, err := c.GetVirtualMachine(vmName)
	if err != nil {
		return err
	}

	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/vms/%s/migrate", strconv.Itoa(vm.ID))
	if err != nil {
		return
	}

	params := map[string]string{}
	if newHead != "" {
		params["new_head_spec"] = newHead
	}

	_, _, err = r.MarshalAndRun(params, nil)
	return
}

func (c *bytemarkClient) ReapVMs() (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/reap_vms")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) DeleteVLAN(id int) (err error) {
	r, err := c.BuildRequest("DELETE", BrainEndpoint, "/admin/vlans/%s", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}
