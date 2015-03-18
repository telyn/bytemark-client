package lib

import (
	"encoding/json"
	"fmt"
)

func (bigv *Client) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	vm = new(VirtualMachine)
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s?view=overview", name.Account, name.Group, name.VirtualMachine)
	data, err := bigv.Request("GET", path, "")

	fmt.Printf("'%s'\r\n", data)

	if err != nil {
		//TODO(telyn): good error handling here
		panic("Couldn't make request")
	}

	err = json.Unmarshal(data, vm)
	if err != nil {
		fmt.Printf("Data returned was not a VirtualMachine\r\n")
		fmt.Printf("%+v\r\n")

		return nil, err
	}
	return vm, nil
}
