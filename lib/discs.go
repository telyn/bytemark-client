package lib

import (
	"encoding/json"
	"fmt"
)

func labelDiscs(discs []*Disc) {
	for i, disc := range discs {
		if disc.Label == "" {
			disc.Label = fmt.Sprintf("%c", 'a'+i)
		}
	}

}

func generateDiscLabel(discIdx int) string {
	return fmt.Sprintf("vd%c", 'a'+discIdx)
}

func (bigv *bigvClient) CreateDiscs(vm VirtualMachineName, discs []*Disc) (err error) {
	err = bigv.validateVirtualMachineName(&vm)
	if err != nil {
		return err
	}
	labelDiscs(discs)

	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/discs", vm.Account, vm.Group, vm.VirtualMachine)
	js, err := json.Marshal(discs)
	if err != nil {
		return err
	}

	_, _, err = bigv.Request(true, "POST", path, string(js))
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
