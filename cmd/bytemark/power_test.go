package main

import (
	"bytemark.co.uk/client/lib"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestResetCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("ResetVirtualMachine", &vmn).Times(1)

	global.App.Run(strings.Split("bytemark reset test-server.test-group.test-account", " "))
	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestRestartCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("RestartVirtualMachine", &vmn).Times(1)

	global.App.Run(strings.Split("bytemark restart test-server.test-group.test-account", " "))
	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestShutdownCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("ShutdownVirtualMachine", &vmn, true).Times(1)

	global.App.Run(strings.Split("bytemark shutdown test-server.test-group.test-account", " "))
	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStartCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("StartVirtualMachine", &vmn).Times(1)

	global.App.Run(strings.Split("bytemark start test-server.test-group.test-account", " "))
	is.Nil(global.Error)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStopCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("StopVirtualMachine", &vmn).Times(1)

	global.App.Run(strings.Split("bytemark stop test-server.test-group.test-account", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
