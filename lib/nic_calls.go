package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func (c *bytemarkClient) AddIP(name *VirtualMachineName, spec *brain.IPCreateRequest) (brain.IPs, error) {
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return nil, err
	}
	nicid := vm.NetworkInterfaces[0].ID

	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/nics/%s/ip_create", name.Account, name.Group, name.VirtualMachine, string(nicid))
	if err != nil {
		return nil, err
	}

	var newSpec *brain.IPCreateRequest
	_, _, err = r.MarshalAndRun(spec, newSpec)
	if err != nil {
		return nil, err
	}
	return newSpec.IPs, nil
}
