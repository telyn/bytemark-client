package lib

import (
	"fmt"
)

func (bigv *Client) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	vm = new(VirtualMachine)
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s?view=overview", name.Account, name.Group, name.VirtualMachine)

	err = bigv.RequestAndUnmarshal("GET", path, "", vm)
	if err != nil {
		return nil, err
	}
	return vm, nil
}
