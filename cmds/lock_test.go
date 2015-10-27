package cmds

import (
	bigv "bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestLockHWProfileCommand(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	vmname := bigv.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account"}
	args := []string{"test-vm.test-group.test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return(args)

	c.When("ParseVirtualMachineName", args[0]).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfileLock", vmname, true).Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfileLock", vmname, false).Return(nil).Times(0)

	cmds := NewCommandSet(config, c)
	cmds.LockHWProfile(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUnlockHWProfileCommand(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	vmname := bigv.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account"}
	args := []string{"test-vm.test-group.test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return(args)

	c.When("ParseVirtualMachineName", args[0]).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfileLock", vmname, true).Return(nil).Times(0)
	c.When("SetVirtualMachineHardwareProfileLock", vmname, false).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.UnlockHWProfile(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
