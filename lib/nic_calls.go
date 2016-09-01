package lib

import (
	"bytes"
	"encoding/json"
)

func (c *bytemarkClient) AddIP(name *VirtualMachineName, spec *IPCreateRequest) (IPs, error) {
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return nil, err
	}
	nicid := vm.NetworkInterfaces[0].ID

	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/nics/%s/ip_create", name.Account, name.Group, name.VirtualMachine, string(nicid))
	if err != nil {
		return nil, err
	}

	js, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	var newSpec *IPCreateRequest
	_, _, err = r.Run(bytes.NewBuffer(js), newSpec)
	if err != nil {
		return nil, err
	}
	return newSpec.IPs, nil
}
