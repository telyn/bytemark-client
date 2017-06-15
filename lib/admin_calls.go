package lib

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"strconv"
)

// UpdateHead is a struct with all the possible settings that can be updated on a head
type UpdateHead struct {
	UsageStrategy   *string
	OvercommitRatio *int
	Label           *string
}

// UpdateTail is a struct with all the possible settings that can be updated on a tail
type UpdateTail struct {
	UsageStrategy   *string
	OvercommitRatio *int
	Label           *string
}

// UpdateStoragePool is a struct with all the possible settings that can be updated on a storage pool
type UpdateStoragePool struct {
	UsageStrategy   *string
	OvercommitRatio *int
	Label           *string
}

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

func (c *bytemarkClient) GetIPRange(idOrCIDR string) (*brain.IPRange, error) {
	if _, err := strconv.Atoi(idOrCIDR); err == nil {
		// Numeric means it is just an ID
		r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/ip_ranges/%s", idOrCIDR)
		if err != nil {
			return nil, err
		}

		var ipRange *brain.IPRange
		_, _, err = r.Run(nil, &ipRange)
		return ipRange, err
	}

	// Non numeric means we got a CIDR
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/ip_ranges?cidr=%s", idOrCIDR)
	if err != nil {
		return nil, err
	}

	// The /admin/ip_ranges endpoint always returns an array of IP ranges,
	// so we just need to get the first one and return it
	var ipRanges []*brain.IPRange
	_, _, err = r.Run(nil, &ipRanges)
	if err != nil {
		return nil, err
	}

	if len(ipRanges) == 0 {
		return nil, fmt.Errorf("IP Range not found")
	}

	if len(ipRanges) > 1 {
		return nil, fmt.Errorf("More than one IP Range found, please report this as a bug")
	}

	return ipRanges[0], nil
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

func (c *bytemarkClient) MigrateVirtualMachine(vmName VirtualMachineName, newHead string) (err error) {
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

func (c *bytemarkClient) AdminCreateGroup(name GroupName, vlanNum int) (err error) {
	err = c.validateGroupName(&name)
	if err != nil {
		return
	}

	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/groups")
	if err != nil {
		return
	}

	obj := map[string]interface{}{
		"account_spec": name.Account,
		"group_name":   name.Group,
	}

	if vlanNum != 0 {
		obj["vlan_num"] = vlanNum
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) CreateIPRange(ipRange string, vlanNum int) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/ip_ranges")
	if err != nil {
		return
	}

	obj := map[string]interface{}{
		"ip_range": ipRange,
		"vlan_num": vlanNum,
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) CancelDiscMigration(id int) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/discs/%s/cancel_migration", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) CancelVMMigration(id int) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/vms/%s/cancel_migration", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) EmptyStoragePool(idOrLabel string) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/storage_pools/%s/empty", idOrLabel)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) EmptyHead(idOrLabel string) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/heads/%s/empty", idOrLabel)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) ReifyDisc(id int) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/discs/%s/reify", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) ApproveVM(name VirtualMachineName, powerOn bool) (err error) {
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return err
	}

	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/vms/%s/approve", strconv.Itoa(vm.ID))
	if err != nil {
		return
	}

	obj := map[string]bool{}
	if powerOn {
		obj["power_on"] = powerOn
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) RejectVM(name VirtualMachineName, reason string) (err error) {
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return err
	}

	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/vms/%s/reject", strconv.Itoa(vm.ID))
	if err != nil {
		return
	}

	obj := map[string]string{
		"reason": reason,
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) RegradeDisc(disc int, newGrade string) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/discs/%s/regrade", strconv.Itoa(disc))
	if err != nil {
		return
	}

	obj := map[string]string{
		"new_grade": newGrade,
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) UpdateVMMigration(name VirtualMachineName, speed *int64, downtime *int) (err error) {
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return err
	}

	r, err := c.BuildRequest("PUT", BrainEndpoint, "/admin/vms/%s/migrate", strconv.Itoa(vm.ID))
	if err != nil {
		return
	}

	obj := map[string]interface{}{}
	if speed != nil {
		obj["migration_speed"] = *speed
	}
	if downtime != nil {
		obj["migration_downtime"] = *downtime
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) CreateUser(username string, privilege string) (err error) {
	r, err := c.BuildRequest("POST", BrainEndpoint, "/admin/users")
	if err != nil {
		return
	}

	obj := map[string]string{
		"username":  username,
		"priv_spec": privilege,
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) UpdateHead(idOrLabel string, options UpdateHead) (err error) {
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/admin/heads/%s", idOrLabel)
	if err != nil {
		return
	}

	obj := map[string]interface{}{}

	if options.OvercommitRatio != nil {
		obj["overcommit_ratio"] = *options.OvercommitRatio
	}
	if options.Label != nil {
		obj["label"] = *options.Label
	}
	if options.UsageStrategy != nil {
		// It is set, but we need to translate an empty string to nil
		if *options.UsageStrategy == "" {
			obj["usage_strategy"] = nil
		} else {
			obj["usage_strategy"] = *options.UsageStrategy
		}
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) UpdateTail(idOrLabel string, options UpdateTail) (err error) {
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/admin/tails/%s", idOrLabel)
	if err != nil {
		return
	}

	obj := map[string]interface{}{}

	if options.OvercommitRatio != nil {
		obj["overcommit_ratio"] = *options.OvercommitRatio
	}
	if options.Label != nil {
		obj["label"] = *options.Label
	}
	if options.UsageStrategy != nil {
		// It is set, but we need to translate an empty string to nil
		if *options.UsageStrategy == "" {
			obj["usage_strategy"] = nil
		} else {
			obj["usage_strategy"] = *options.UsageStrategy
		}
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}

func (c *bytemarkClient) UpdateStoragePool(idOrLabel string, options UpdateStoragePool) (err error) {
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/admin/storage_pools/%s", idOrLabel)
	if err != nil {
		return
	}

	obj := map[string]interface{}{}

	if options.OvercommitRatio != nil {
		obj["overcommit_ratio"] = *options.OvercommitRatio
	}
	if options.Label != nil {
		obj["label"] = *options.Label
	}
	if options.UsageStrategy != nil {
		// It is set, but we need to translate an empty string to nil
		if *options.UsageStrategy == "" {
			obj["usage_strategy"] = nil
		} else {
			obj["usage_strategy"] = *options.UsageStrategy
		}
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return
}
