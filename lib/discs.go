package lib

import (
	"bytes"
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
		return
	}
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return
	}
	discs := []Disc{disc}
	labelDiscs(discs, len(vm.Discs))

	r, err := c.BuildRequest("POST", EP_BRAIN, "/accounts/%s/groups/%s/virtual_machines/%s/discs", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return
	}

	js, err := json.Marshal(discs[0])
	if err != nil {
		return
	}

	_, _, err = r.Run(bytes.NewBuffer(js), nil)
	return

}

func (c *bytemarkClient) DeleteDisc(vm VirtualMachineName, discLabelOrID string) (err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("DELETE", EP_BRAIN, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s?purge=true", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)

	return
}

func (c *bytemarkClient) ResizeDisc(vm VirtualMachineName, discLabelOrID string, sizeMB int) (err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", EP_BRAIN, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	// TODO(telyn): marshal json instead of sprintf
	disc := fmt.Sprintf(`{"size":%d}`, sizeMB)

	_, _, err = r.Run(bytes.NewBufferString(disc), nil)
	return err
}

func (c *bytemarkClient) GetDisc(vm VirtualMachineName, discLabelOrID string) (disc *Disc, err error) {
	err = c.validateVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("GET", EP_BRAIN, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, disc)
	return
}
