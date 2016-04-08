package main

import (
	"bytemark.co.uk/client/lib"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestLockHWProfileCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfileLock", &vmname, true).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark lock hwprofile test-server.test-group.test-account", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUnlockHWProfileCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfileLock", &vmname, false).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark unlock hwprofile test-server.test-group.test-account", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
