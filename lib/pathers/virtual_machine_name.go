package pathers

import (
	"errors"
	"fmt"
)

// VirtualMachineName is the triplet-form of the name of a VirtualMachine, which should be enough to find the VM.
type VirtualMachineName struct {
	GroupName
	VirtualMachine string
}

func (vm VirtualMachineName) String() string {
	vm.defaultIfNeeded()
	if vm.Account == "" {
		return fmt.Sprintf("%s.%s", vm.VirtualMachine, vm.Group)
	}
	return fmt.Sprintf("%s.%s.%s", vm.VirtualMachine, vm.Group, vm.Account)
}

// VirtualMachinePath returns a Brain URL path for this virtual machine, or an
// error if one cannot be made.
func (vm VirtualMachineName) VirtualMachinePath() (string, error) {
	vm.defaultIfNeeded()
	if vm.VirtualMachine == "" {
		return "", errors.New("VirtualMachine component of VirtualMachineName is empty")
	}
	groupPath, err := vm.GroupPath()
	return groupPath + fmt.Sprintf("/virtual_machines/%s", vm.VirtualMachine), err
}
