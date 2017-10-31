package main

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestSetCDROM(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("SetVirtualMachineCDROM", vmname, "test-cdrom").Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark set cdrom test-server.test-group.test-account test-cdrom", " "))
	is.Nil(err)
	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}
}

func TestSetCores(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	vm := getFixtureVM()
	c.When("GetVirtualMachine", vmname).Return(&vm)
	c.When("SetVirtualMachineCores", vmname, 4).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark set cores --force test-server.test-group.test-account 4", " "))
	is.Nil(err)

	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}
}

func TestSetMemory(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account"}

	config.When("GetVirtualMachine").Return(defVM)

	vm := getFixtureVM()
	c.When("GetVirtualMachine", vmname).Return(vm)
	c.When("SetVirtualMachineMemory", vmname, 4096).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark set memory --force test-server 4", " "))
	is.Nil(err)

	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}

	config, c, app = testutil.BaseTestAuthSetup(t, false, commands)
	config.When("GetVirtualMachine").Return(defVM)

	c.When("GetVirtualMachine", vmname).Return(vm)
	c.When("SetVirtualMachineMemory", vmname, 16384).Return(nil).Times(1)

	err = app.Run(strings.Split("bytemark set memory --force test-server 16384M", " "))
	if err != nil {
		t.Error(err)
	}

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestSetHWProfileCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestSetup(t, false, commands)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(defVM)

	// test no arguments, nothing should happen
	c.When("SetVirtualMachineHardwareProfile", vmname).Return(nil).Times(0) // don't do anything
	c.When("AuthWithToken", "test-token").Return(nil).Times(0)

	err := app.Run(strings.Split("bytemark set hwprofile test-server", " "))
	is.NotNil(err) // TODO(telyn): actually check error type

	if ok, vErr := c.Verify(); !ok {
		t.Error(vErr)
	}

	// test hardware profile only

	config, c, app = testutil.BaseTestAuthSetup(t, false, commands)
	config.When("GetVirtualMachine").Return(defVM)

	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool(nil)).Return(nil).Times(1)

	err = app.Run(strings.Split("bytemark set hwprofile test-server virtio123", " "))
	is.Nil(err)

	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}

	// test --lock flag
	config, c, app = testutil.BaseTestAuthSetup(t, false, commands)
	config.When("GetVirtualMachine").Return(defVM)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool{true}).Return(nil).Times(1)

	err = app.Run(strings.Split("bytemark set hwprofile --lock test-server virtio123", " "))
	is.Nil(err)

	if ok, vErr := c.Verify(); !ok {
		t.Fatal(vErr)
	}

	// test --unlock flag
	config, c, app = testutil.BaseTestAuthSetup(t, false, commands)
	config.When("GetVirtualMachine").Return(defVM)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool{false}).Return(nil).Times(1)

	err = app.Run(strings.Split("bytemark set hwprofile --unlock test-server virtio123", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
