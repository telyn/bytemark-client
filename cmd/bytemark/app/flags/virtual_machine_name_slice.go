package flags

// This file was automatically generate using
// cmd/bytemark/app/flags/gen/slice_flags - do not edit it by hand!

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
)

// VirtualMachineNameSliceFlag is used for VirtualMachineNameFlags that may be specified more than
// once. It's a slice of VirtualMachineNameFlag in order to avoid rewriting parsing
// logic.
type VirtualMachineNameSliceFlag []VirtualMachineNameFlag

// Preprocess calls Preprocess on all the underlying VirtualMachineNameFlags
func (sf *VirtualMachineNameSliceFlag) Preprocess(ctx *app.Context) error {
	for i := range *sf {
		err := (*sf)[i].Preprocess(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set appends a VirtualMachineNameFlag (created for you) to the slice
func (sf *VirtualMachineNameSliceFlag) Set(value string) error {
	flag := VirtualMachineNameFlag{}
	err := flag.Set(value)
	if err != nil {
		return err
	}
	*sf = append(*sf, flag)
	return nil
}

// String returns all values in the slice, comma-delimeted
func (sf VirtualMachineNameSliceFlag) String() string {
	strs := make([]string, len(sf))
	for i, value := range sf {
		strs[i] = value.String()
	}
	return strings.Join(strs, ", ")
}

// VirtualMachineNameSlice returns the named flag as a VirtualMachineNameSliceFlag,
// if it was one in the first place.
func VirtualMachineNameSlice(ctx *app.Context, name string) VirtualMachineNameSliceFlag {
	if sf, ok := ctx.Context.Generic(name).(*VirtualMachineNameSliceFlag); ok {
		return *sf
	}
	return VirtualMachineNameSliceFlag{}
}
