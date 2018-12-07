package flags

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// VirtualMachineNameFlag is used for all --server flags, or should be at least.
type VirtualMachineNameFlag struct {
	VirtualMachineName *lib.VirtualMachineName
	Value              string
}

// Set runs lib.ParseVirtualMachineName using the c.Client() to make sure we have a valid group name
func (name *VirtualMachineNameFlag) Set(value string) error {
	name.Value = value
	return nil
}

// Preprocess defaults the value of this flag to the default server from the
// config attached to the context and then runs lib.ParseVirtualMachineName
// This is an implementation of `app.Preprocessor`, which is detected and
// called automatically by actions created with `app.Action`
func (name *VirtualMachineNameFlag) Preprocess(c *app.Context) (err error) {
	if name.Value == "" {
		return
	}
	vmName, err := lib.ParseVirtualMachineName(name.Value, c.Config().GetVirtualMachine())
	name.VirtualMachineName = &vmName
	return
}

// String returns the VirtualMachineName as a string.
func (name VirtualMachineNameFlag) String() string {
	if name.VirtualMachineName != nil {
		return name.VirtualMachineName.String()
	}
	return ""
}
