package lib

import (
	"fmt"
)

// VirtualMachineName is the triplet-form of the name of a VirtualMachine, which should be enough to find the VM.
// VirtualMachineName implements the Pather interface
type VirtualMachineName struct {
	VirtualMachine string
	Group          string
	Account        string
}

func (vm VirtualMachineName) String() string {
	if vm.Group == "" {
		vm.Group = DefaultGroup
	}
	if vm.Account == "" {
		return fmt.Sprintf("%s.%s", vm.VirtualMachine, vm.Group)
	}
	return fmt.Sprintf("%s.%s.%s", vm.VirtualMachine, vm.Group, vm.Account)
}

// GroupName returns the group and account of this VirtualMachineName as a group.
func (vm VirtualMachineName) GroupName() GroupName {
	return GroupName{
		Group:   vm.Group,
		Account: vm.Account,
	}
}

// Path returns the URL path for this VM, if possible.
// If the VM is not full specified (i.e. does not have an account, group and
// name), it instead returns an error.
func (vm VirtualMachineName) VirtualMachinePath() (string, error) {
	if vm.VirtualMachine == "" || vm.Group == "" || vm.Account == "" {
		return "", fmt.Errorf("Server %q was not fully specified so cannot make a URL", vm)
	}
	path := fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s",
		vm.Account, vm.Group, vm.VirtualMachine)
	return path, nil
}
