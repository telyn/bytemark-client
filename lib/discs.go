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

func (disc *Disc) Validate() (*Disc, error) {
	if disc.StorageGrade == "" {
		newDisc := *disc
		newDisc.StorageGrade = "sata"
		return &newDisc, nil
	}
	return disc, nil
}

func (c *bytemarkClient) CreateDisc(name VirtualMachineName, disc Disc) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return err
	}
	discs := []Disc{disc}
	labelDiscs(discs, len(vm.Discs))

	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/discs", name.Account, name.Group, name.VirtualMachine)
	js, err := json.Marshal(discs[0])
	if err != nil {
		return err
	}

	_, err = c.RequestAndRead(true, "POST", path, string(js))
	return err

}

func (c *bytemarkClient) DeleteDisc(vm VirtualMachineName, discLabelOrID string) (err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/discs/%s?purge=true", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)

	_, _, err = c.Request(true, "DELETE", path, "")

	return err
}

func (c *bytemarkClient) ResizeDisc(vm VirtualMachineName, discLabelOrID string, sizeMB int) (err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/discs/%s", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)

	disc := fmt.Sprintf(`{"size":%d}`, sizeMB)

	_, _, err = c.Request(true, "PUT", path, disc)
	return err
}

func (c *bytemarkClient) GetDisc(vm VirtualMachineName, discLabelOrID string) (disc *Disc, err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return nil, err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/discs/%s", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)

	disc = new(Disc)

	err = c.RequestAndUnmarshal(true, "GET", path, "", disc)
	return disc, err
}
