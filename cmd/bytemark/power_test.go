package main

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func TestResetCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("ResetVirtualMachine", vmn).Times(1)

	err := app.Run(strings.Split("bytemark reset test-server.test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestRestartCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("ShutdownVirtualMachine", vmn, true).Times(1)
	c.When("GetVirtualMachine", vmn).Return(brain.VirtualMachine{PowerOn: false})
	c.When("StartVirtualMachine", vmn).Times(1)

	err := app.Run(strings.Split("bytemark restart test-server.test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestShutdownCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("ShutdownVirtualMachine", vmn, true).Times(1)
	c.When("GetVirtualMachine", vmn).Return(brain.VirtualMachine{PowerOn: false})

	err := app.Run(strings.Split("bytemark shutdown test-server.test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStartCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("StartVirtualMachine", vmn).Times(1)

	err := app.Run(strings.Split("bytemark start test-server.test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStopCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("StopVirtualMachine", vmn).Times(1)

	err := app.Run(strings.Split("bytemark stop test-server.test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
