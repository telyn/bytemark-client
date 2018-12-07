package flags

import (
	"flag"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// VirtualMachineNameSliceFlag is used for --server flags that may be specified more than
// once. It's a slice of VirtualMachineNameFlag in order to avoid rewriting parsing
// logic.
type VirtualMachineNameSliceFlag interface {
	flag.Value
	app.Preprocesser
	VirtualMachineNames() []lib.VirtualMachineName
}

type virtualMachineNameSliceFlag struct {
	GenericSliceFlag
}

// NewVirtualMachineNameSliceFlag creates a new
func NewVirtualMachineNameSliceFlag() VirtualMachineNameSliceFlag {
	return virtualMachineNameSliceFlag{
		template: VirtualMachineNameFlag{},
	}
}

// VirtualMachineNames returns the VirtualMachineNameSlice for which this flag is named
func (gnsf VirtualMachineNameSliceFlag) VirtualMachineNames() (virtualMachineNames []lib.VirtualMachineName) {
	gnsf.copyValues(virtualMachineNames)
	return
}
