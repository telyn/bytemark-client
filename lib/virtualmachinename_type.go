package lib

import (
	"fmt"
)

// VirtualMachineName is the triplet-form of the name of a VirtualMachine, which should be enough to find the VM.
type VirtualMachineName struct {
	VirtualMachine string
	Group          string
	Account        string
}

func (vm VirtualMachineName) String() string {
	if vm.Group == "" {
		vm.Group = "default"
	}
	if vm.Account == "" {
		return fmt.Sprintf("%s.%s", vm.VirtualMachine, vm.Group)
	}
	return fmt.Sprintf("%s.%s.%s", vm.VirtualMachine, vm.Group, vm.Account)
}

func (vm VirtualMachineName) GroupName() *GroupName {
	return &GroupName{
		Group:   vm.Group,
		Account: vm.Account,
	}
}
