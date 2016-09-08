package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestSetCores(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	vm := getFixtureVM()
	c.When("GetVirtualMachine", &vmname).Return(&vm)
	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineCores", &vmname, 4).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark set cores --force test-server.test-group.test-account 4", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestSetMemory(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	vm := getFixtureVM()
	c.When("GetVirtualMachine", &vmname).Return(&vm)
	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineMemory", &vmname, 4096).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark set memory --force test-server 4", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	config.Reset()
	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.Reset()
	c.When("GetVirtualMachine", &vmname).Return(&vm)
	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineMemory", &vmname, 16384).Return(nil).Times(1)

	err = global.App.Run(strings.Split("bytemark set memory --force test-server 16384M", " "))
	if err != nil {
		t.Error(err)
	}

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestSetHWProfileCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	// test no arguments, nothing should happen
	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vmname)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", &vmname).Return(nil).Times(0) // don't do anything

	err := global.App.Run(strings.Split("bytemark set hwprofile test-server", " "))
	is.NotNil(err) // TODO(telyn): actually check error type

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test hardware profile only

	config.Reset()
	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	c.Reset()
	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", &vmname, "virtio123", []bool(nil)).Return(nil).Times(1)

	err = global.App.Run(strings.Split("bytemark set hwprofile test-server virtio123", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test --lock flag
	c.Reset()
	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", &vmname, "virtio123", []bool{true}).Return(nil).Times(1)

	err = global.App.Run(strings.Split("bytemark set hwprofile --lock test-server virtio123", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test --unlock flag
	c.Reset()
	c.When("ParseVirtualMachineName", "test-server", []*lib.VirtualMachineName{&defVM}).Return(&vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", &vmname, "virtio123", []bool{false}).Return(nil).Times(1)

	err = global.App.Run(strings.Split("bytemark set hwprofile --unlock test-server virtio123", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
