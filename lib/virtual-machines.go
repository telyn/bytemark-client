package lib

import (
	"fmt"
)

// GetVirtualMachine
func (bigv *BigVClient) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	vm = new(VirtualMachine)
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s?include_deleted=true&view=overview", name.Account, name.Group, name.VirtualMachine)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", vm)
	return vm, err
}

// returns nil on success. Probably.
func (bigv *BigVClient) DeleteVirtualMachine(name VirtualMachineName, purge bool) (err error) {
	// TODO(telyn): URL escaping
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if purge {
		path += "?purge=true"
	}

	_, _, err = bigv.Request(true, "DELETE", path, "")
	return err
}

func (bigv *BigVClient) UndeleteVirtualMachine(name VirtualMachineName) (err error) {
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"deleted":false}`)
	return err
}
