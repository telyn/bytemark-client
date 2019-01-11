package lib

import (
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func (c *bytemarkClient) AddIP(name VirtualMachineName, spec brain.IPCreateRequest) (brain.IPs, error) {
	err := c.checkVirtualMachinePather(&name)
	if err != nil {
		return nil, err
	}
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return nil, err
	}
	nicid := strconv.Itoa(vm.NetworkInterfaces[0].ID)

	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/nics/%s/ip_create", name.Account, name.Group, name.VirtualMachine, nicid)
	if err != nil {
		return nil, err
	}

	var createdSpec brain.IPCreateRequest
	_, _, err = r.MarshalAndRun(spec, &createdSpec)
	if err != nil {
		return nil, err
	}
	return createdSpec.IPs, nil
}
