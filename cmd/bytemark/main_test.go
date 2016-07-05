package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

var defVM lib.VirtualMachineName
var defGroup lib.GroupName

func baseTestSetup() (config *mocks.Config, client *mocks.Client) {
	config = new(mocks.Config)
	client = new(mocks.Client)
	global.Client = client
	global.Config = config

	baseAppSetup()
	return
}

func traverseAllCommands(cmds []cli.Command, fn func(cli.Command)) {
	if cmds == nil {
		return
	}
	for _, c := range cmds {
		fn(c)
		traverseAllCommands(c.Subcommands, fn)
	}
}

func getFixtureVM() lib.VirtualMachine {
	return lib.VirtualMachine{
		Name:     "test-server",
		Hostname: "test-server.test-group",
		GroupID:  1,
		ZoneName: "test-zone",
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
