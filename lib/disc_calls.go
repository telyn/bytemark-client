package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func labelDiscs(discs []brain.Disc, offset ...int) {
	realOffset := 0
	if len(offset) >= 1 {
		realOffset = offset[0]
	}
	for i := range discs {
		if discs[i].Label == "" {
			discs[i].Label = fmt.Sprintf("disc-%d", realOffset+i+1)
		}
	}

}

// CreateDisc creates the given Disc and attaches it to the given virtual machine.
func (c *bytemarkClient) CreateDisc(name *VirtualMachineName, disc brain.Disc) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return
	}
	vm, err := c.GetVirtualMachine(name)
	if err != nil {
		return
	}
	discs := []brain.Disc{disc}

	labelDiscs(discs, vm.GetDiscLabelOffset())

	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs", name.Account, name.Group, name.VirtualMachine)
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

// DeleteDisc removes the specified disc from the given virtual machine
func (c *bytemarkClient) DeleteDisc(vm *VirtualMachineName, discLabelOrID string) (err error) {
	err = c.validateVirtualMachineName(vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("DELETE", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s?purge=true", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)

	return
}

// ResizeDisc resizes the specified disc to the given size in megabytes
func (c *bytemarkClient) ResizeDisc(vm *VirtualMachineName, discLabelOrID string, sizeMB int) (err error) {
	err = c.validateVirtualMachineName(vm)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	// TODO(telyn): marshal json instead of sprintf
	disc := fmt.Sprintf(`{"size":%d}`, sizeMB)

	_, _, err = r.Run(bytes.NewBufferString(disc), nil)
	return err
}

// GetDisc returns the specified disc from the given virtual machine.
func (c *bytemarkClient) GetDisc(vm *VirtualMachineName, discLabelOrID string) (disc *brain.Disc, err error) {
	disc = new(brain.Disc)
	err = c.validateVirtualMachineName(vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)

	if err != nil {
		return
	}

	_, _, err = r.Run(nil, disc)
	return
}
