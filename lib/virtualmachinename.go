package lib

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// VirtualMachineName is the triplet-form of the name of a VirtualMachine, which should be enough to find the VM.
type VirtualMachineName struct {
	VirtualMachine string
	Group          string
	Account        pathers.AccountName
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
