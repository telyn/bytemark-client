package lib

import "encoding/json"

// GetVirtualMachine requests an overview of the named VM, regardless of its deletion status.
func (bigv *bigvClient) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	vm = new(VirtualMachine)
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s?include_deleted=true&view=overview", name.Account, name.Group, name.VirtualMachine)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", vm)
	if err != nil {
		vm = nil
	}
	return vm, err
}

// DeleteVirtualMachine deletes the named virtual machine.
// returns nil on success or an error otherwise.
func (bigv *bigvClient) DeleteVirtualMachine(name VirtualMachineName, purge bool) (err error) {
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if purge {
		path += "?purge=true"
	}

	_, _, err = bigv.Request(true, "DELETE", path, "")
	return err
}

// UndeleteVirtualMachine changes the deleted flag on a VM back to false.
// Return nil on success, an error otherwise.
func (bigv *bigvClient) UndeleteVirtualMachine(name VirtualMachineName) (err error) {
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"deleted":false}`)
	return err
}

//CreateVirtualMachine creates a virtual machine in the given group.
func (bigv *bigvClient) CreateVirtualMachine(group GroupName, spec VirtualMachineSpec) (vm *VirtualMachine, err error) {
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/create_vm", group.Account, group.Group)
	// TODO(telyn): possibly better to build a map[string]string and marshal that
	js, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	vm = new(VirtualMachine)
	err = bigv.RequestAndUnmarshal(true, "POST", path, string(js), vm)
	return vm, err
}
