package lib

import (
	"encoding/json"
	"fmt"
)

func labelDiscs(discs []Disc, offset ...int) {
	realOffset := 0
	if len(offset) >= 1 {
		realOffset = offset[0]
	}
	for i := range discs {
		if discs[i].Label == "" {
			discs[i].Label = fmt.Sprintf("vd%c", 'a'+realOffset+i)
		}
	}

}

func (bigv *bigvClient) CreateDisc(name VirtualMachineName, disc Disc) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	vm, err := bigv.GetVirtualMachine(name)
	// smell bad
	discs := []Disc{disc}
	labelDiscs(discs, len(vm.Discs))

	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/discs", name.Account, name.Group, name.VirtualMachine)
	js, err := json.Marshal(discs[0])
	if err != nil {
		return err
	}

	_, err = bigv.RequestAndRead(true, "POST", path, string(js))
	return err

}

func (bigv *bigvClient) DeleteDisc(vm VirtualMachineName, discID int) (err error) {
	err = bigv.validateVirtualMachineName(&vm)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/discs/%s?purge=true", vm.Account, vm.Group, vm.VirtualMachine, fmt.Sprintf("%d", discID))

	_, _, err = bigv.Request(true, "DELETE", path, "")

	return err
}
