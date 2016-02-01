package cmds

import (
	bigv "bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"github.com/cheekybits/is"
	"testing"
)

func TestSetCores(t *testing.T) {
	is := is.New(t)
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	vmname := bigv.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account"}
	args := []string{"test-vm.test-group.test-account", "4"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return(args)
	config.When("GetVirtualMachine").Return(bigv.VirtualMachineName{})

	c.When("ParseVirtualMachineName", args[0], []bigv.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineCores", vmname, 4).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	is.Equal(0, cmds.SetCores(args))

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestSetMemory(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	vmname := bigv.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account"}
	args := []string{"test-vm.test-group.test-account", "4"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return(args)
	config.When("GetVirtualMachine").Return(bigv.VirtualMachineName{})

	c.When("ParseVirtualMachineName", args[0], []bigv.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineMemory", vmname, 4096).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.SetMemory(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	args = []string{"test-vm.test-group.test-account", "16384M"}

	config.Reset()
	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return(args)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(bigv.VirtualMachineName{})

	c.Reset()
	c.When("ParseVirtualMachineName", args[0], []bigv.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineMemory", vmname, 16384).Return(nil).Times(1)

	cmds = NewCommandSet(config, c)
	cmds.SetMemory(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestSetHWProfileCommand(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	vmname := bigv.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account"}
	args := []string{"test-vm.test-group.test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return(args)
	config.When("GetVirtualMachine").Return(bigv.VirtualMachineName{})

	// test no arguments, nothing should happen
	c.When("ParseVirtualMachineName", args[0], []bigv.VirtualMachineName{{}}).Return(vmname)
	c.When("AuthWithToken", "test-token").Return(nil).Times(0)              // don't talk to BigV
	c.When("SetVirtualMachineHardwareProfile", vmname).Return(nil).Times(0) // don't do anything

	cmds := NewCommandSet(config, c)
	cmds.SetHWProfile(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test hardware profile only
	args = []string{"test-vm.test-group.test-account", "virtio123"}

	config.Reset()
	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return(args)
	config.When("GetVirtualMachine").Return(bigv.VirtualMachineName{})

	c.Reset()
	c.When("ParseVirtualMachineName", args[0], []bigv.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool(nil)).Return(nil).Times(1)

	cmds = NewCommandSet(config, c)
	cmds.SetHWProfile(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test --lock flag
	c.Reset()
	c.When("ParseVirtualMachineName", args[0], []bigv.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool{true}).Return(nil).Times(1)

	args_flag := []string{"--lock", "test-vm.test-group.test-account", "virtio123"}

	cmds = NewCommandSet(config, c)
	cmds.SetHWProfile(args_flag)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test --unlock flag
	c.Reset()
	c.When("ParseVirtualMachineName", args[0], []bigv.VirtualMachineName{{}}).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", []bool{false}).Return(nil).Times(1)

	args_flag = []string{"--unlock", "test-vm.test-group.test-account", "virtio123"}

	cmds = NewCommandSet(config, c)
	cmds.SetHWProfile(args_flag)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
