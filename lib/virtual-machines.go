package lib

// GetVirtualMachine requests an overview of the named VM, regardless of its deletion status.
func (bigv *BigVClient) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	vm = new(VirtualMachine)
	path := BuildUrl("/accounts/%s/groups/%s/virtual_machines/%s?include_deleted=true&view=overview", name.Account, name.Group, name.VirtualMachine)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", vm)
	if err != nil {
		vm = nil
	}
	return vm, err
}

// DeleteVirtualMachine deletes the named virtual machine.
// returns nil on success or an error otherwise.
func (bigv *BigVClient) DeleteVirtualMachine(name VirtualMachineName, purge bool) (err error) {
	// TODO(telyn): URL escaping
	path := BuildUrl("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if purge {
		path += "?purge=true"
	}

	_, _, err = bigv.Request(true, "DELETE", path, "")
	return err
}

// UndeleteVirtualMachine changes the deleted flag on a VM back to false.
// Return nil on success, an error otherwise.
func (bigv *BigVClient) UndeleteVirtualMachine(name VirtualMachineName) (err error) {
	path := BuildUrl("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"deleted":false}`)
	return err
}
