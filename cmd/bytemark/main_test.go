package main

import (
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func baseTestSetup() (config *mocks.Config, client *mocks.Client) {
	config = new(mocks.Config)
	client = new(mocks.Client)
	global.Client = client
	global.Config = config

	baseAppSetup()
	return
}

func getFixtureVM() lib.VirtualMachine {
	return lib.VirtualMachine{
		Name:     "test-server",
		Hostname: "test-server.test-group",
		GroupID:  1,
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
