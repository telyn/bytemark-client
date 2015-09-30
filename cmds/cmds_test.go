package cmds

import (
	bigv "bytemark.co.uk/client/lib"
	//"github.com/cheekybits/is"
)

///////////////////////
// Support Functions //
///////////////////////

func getFixtureVM() bigv.VirtualMachine {
	return bigv.VirtualMachine{
		Name:    "test-vm",
		GroupID: 1,
	}
}

func getFixtureGroup() bigv.Group {
	vms := make([]*bigv.VirtualMachine, 1, 1)
	vm := getFixtureVM()
	vms[0] = &vm

	return bigv.Group{
		Name:            "test-group",
		VirtualMachines: vms,
	}
}

////////////////
// Test Cases //
////////////////
