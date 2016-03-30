package main

import (
	"bytemark.co.uk/client/lib"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestSetCores(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server", []lib.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineCores", vmname, 4).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark set cores test-server.test-group.test-account 4", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestSetMemory(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server", []lib.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineMemory", vmname, 4096).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark set memory test-server 4", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	config.Reset()
	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.Reset()
	c.When("ParseVirtualMachineName", "test-server", []lib.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineMemory", vmname, 16384).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark set memory test-server 16384M", " "))

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestSetHWProfileCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup()

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	// test no arguments, nothing should happen
	c.When("ParseVirtualMachineName", "test-server", []lib.VirtualMachineName{{}}).Return(vmname)
	c.When("AuthWithToken", "test-token").Return(nil).Times(0)              // don't talk to the API
	c.When("SetVirtualMachineHardwareProfile", vmname).Return(nil).Times(0) // don't do anything

	global.App.Run(strings.Split("bytemark set hwprofile test-server", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test hardware profile only

	config.Reset()
	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.Reset()
	c.When("ParseVirtualMachineName", "test-server", []lib.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool(nil)).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark set hwprofile test-server virtio123", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test --lock flag
	c.Reset()
	c.When("ParseVirtualMachineName", "test-server", []lib.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool{true}).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark set hwprofile --lock test-server virtio123", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test --unlock flag
	c.Reset()
	c.When("ParseVirtualMachineName", "test-server", []lib.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool{false}).Return(nil).Times(1)

	global.App.Run(strings.Split("bytemark set hwprofile --unlock test-server virtio123", " "))
	is.Nil(global.Error)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
