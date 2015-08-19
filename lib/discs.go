package lib

import "fmt"

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

func (bigv *bigvClient) DeleteDisc(vm VirtualMachineName, discID int) (err error) {
	err = bigv.validateVirtualMachineName(&vm)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/discs/%s", vm.Account, vm.Group, vm.VirtualMachine, fmt.Sprintf("%d", discID))

	_, _, err = bigv.Request(true, "DELETE", path, "")
	return err
}
