package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	//"github.com/cheekybits/is"
)

///////////////////////
// Support Functions //
///////////////////////

func getFixtureVM() lib.VirtualMachine {
	return lib.VirtualMachine{
		Name:    "test-server",
		GroupID: 1,
	}
}

func getFixtureGroup() lib.Group {
	vms := make([]*lib.VirtualMachine, 1, 1)
	vm := getFixtureVM()
	vms[0] = &vm

	return lib.Group{
		Name:            "test-group",
		VirtualMachines: vms,
	}
}

////////////////
// Test Cases //
////////////////
