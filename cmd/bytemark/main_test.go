package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"testing"
)

var defVM lib.VirtualMachineName
var defGroup lib.GroupName

func baseTestSetup(t *testing.T, admin bool) (config *mocks.Config, client *mocks.Client) {
	config = new(mocks.Config)
	client = new(mocks.Client)
	config.When("GetBool", "admin").Return(admin, nil)
	global.Client = client
	global.Config = config

	app, err := baseAppSetup()
	if err != nil {
		t.Fatal(err)
	}
	global.App = app
	oldWriter := global.App.Writer
	global.App.Writer = ioutil.Discard
	for _, c := range commands {
		//config.When("Get", "token").Return("no-not-a-token")

		// the issue is that Command.FullName() is dependent on Command.commandNamePath.
		// Command.commandNamePath is filled in when the parent's Command.startApp is called
		// and startApp is only called when you actually try to run that command or one of
		// its subcommands. So we run "bytemark <command> help" on all commands that have
		// subcommands in order to get every subcommand to have a correct Command.commandPath

		if c.Subcommands != nil && len(c.Subcommands) > 0 {
			fmt.Fprintf(os.Stderr, c.Name)
			_ = global.App.Run([]string{"bytemark.test", c.Name, "help"})
		}
	}
	global.App.Writer = oldWriter
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

func getFixtureVM() brain.VirtualMachine {
	return brain.VirtualMachine{
		Name:     "test-server",
		Hostname: "test-server.test-group",
		GroupID:  1,
	}
}

func getFixtureGroup() brain.Group {
	vms := make([]*brain.VirtualMachine, 1, 1)
	vm := getFixtureVM()
	vms[0] = &vm

	return brain.Group{
		Name:            "test-group",
		VirtualMachines: vms,
	}
}
