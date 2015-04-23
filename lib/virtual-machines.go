package lib

import (
	"fmt"
)

// GetVirtualMachine
func (bigv *BigVClient) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	vm = new(VirtualMachine)
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s?view=overview", name.Account, name.Group, name.VirtualMachine)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", vm)
	if err != nil {
		return nil, bigv.PopulateError(err, name.String(), "virtual machine", "delete")
	}
	return vm, nil
}

// returns nil on success. Probably.
func (bigv *BigVClient) DeleteVirtualMachine(name VirtualMachineName, purge bool) (err error) {
	// TODO(telyn): URL escaping
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if purge {
		path += "?purge=true"
	}

	_, _, err = bigv.Request(true, "DELETE", path, "")
	if err != nil {
		return bigv.PopulateError(err, name.String(), "virtual machine", "delete")
	}
	return nil
}

func (bigv *BigVClient) UndeleteVirtualMachine(name VirtualMachineName) (err error) {
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, "{deleted:false}")
	if err != nil {
		return bigv.PopulateError(err, name.String(), "virtual machine", "delete")
	}
	return nil
}
