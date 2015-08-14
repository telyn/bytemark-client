package main

import (
	bigv "bigv.io/client/lib"
	"testing"
	//"github.com/cheekybits/is"
)

///////////////////////
// Support Functions //
///////////////////////

func getFixtureVM() bigv.VirtualMachine {
	return bigv.VirtualMachine{
		Name:    "test-vm",
		GroupID: 1,
	}
}

func getFixtureGroup() bigv.Group {
	vms := make([]*bigv.VirtualMachine, 1, 1)
	vm := getFixtureVM()
	vms[0] = &vm

	return bigv.Group{
		Name:            "test-group",
		VirtualMachines: vms,
	}
}

////////////////
// Test Cases //
////////////////

func TestCommandConfig(t *testing.T) {
	config := &mockConfig{}

	config.When("GetV", "user").Return(ConfigVar{"user", "old-test-user", "config"})
	config.When("Get", "user").Return("old-test-user")
	config.When("Silent").Return(true)

	config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

	cmds := NewCommandSet(config, nil)
	cmds.Config([]string{"set", "user", "test-user"})

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateVMCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm"})

	c.When("ParseVirtualMachineName", "test-vm").Return(bigv.VirtualMachineName{VirtualMachine: "test-vm"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	vm := bigv.VirtualMachineSpec{
		Discs: []*bigv.Disc{
			&bigv.Disc{
				Size:         25 * 1024,
				StorageGrade: "sata",
			},
		},
		VirtualMachine: &bigv.VirtualMachine{
			Name:                  "test-vm",
			Autoreboot:            true,
			Cores:                 1,
			Memory:                1,
			CdromURL:              "https://example.com/example.iso",
			HardwareProfile:       "test-profile",
			HardwareProfileLocked: true,
			ZoneName:              "test-zone",
		},
		Reimage: &bigv.ImageInstall{
			Distribution: "test-image",
			RootPassword: "test-password",
		},
	}

	group := bigv.GroupName{
		Group:   "test-group",
		Account: "test-account",
	}

	c.When("CreateVirtualMachine", group, vm).Return(vm.VirtualMachine, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.CreateVM([]string{
		"--account", "test-account",
		"--cdrom", "https://example.com/example.iso",
		"--cores", "1",
		"--discs", "25",
		"--group", "test-group",
		"--hwprofile", "test-profile",
		"--hwprofile-locked",
		"--image", "test-image",
		"--memory", "1",
		"--root-password", "test-password",
		"--zone", "test-zone",
		"test-vm",
	})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowGroupCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-group.test-account"})

	c.When("ParseGroupName", "test-group.test-account").Return(bigv.GroupName{Group: "test-group", Account: "test-account"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	group := getFixtureGroup()
	c.When("GetGroup", bigv.GroupName{Group: "test-group", Account: "test-account"}).Return(&group, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ShowGroup([]string{"test-group.test-account"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestResetCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}
	vmn := bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm.test-group.test-account"})
	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account").Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("ResetVirtualMachine", vmn).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ResetVM([]string{"test-vm.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestRestartCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}
	vmn := bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm.test-group.test-account"})
	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account").Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("RestartVirtualMachine", vmn).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.Restart([]string{"test-vm.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestShutdownCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}
	vmn := bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm.test-group.test-account"})
	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account").Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	c.When("ShutdownVirtualMachine", vmn, true).Times(1)
	cmds.Shutdown([]string{"test-vm.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStartCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}
	vmn := bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm.test-group.test-account"})
	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account").Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	cmds := NewCommandSet(config, c)

	c.When("StartVirtualMachine", vmn).Times(1)
	cmds.Start([]string{"test-vm.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStopCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}

	vmn := bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm.test-group.test-account"})
	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account").Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("StopVirtualMachine", vmn).Times(1)

	cmds := NewCommandSet(config, c)

	cmds.Stop([]string{"test-vm.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowVMCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm.test-group.test-account"})

	c.When("ParseVirtualMachineName", "test-vm.test-group.test-account").Return(bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	vm := getFixtureVM()
	c.When("GetVirtualMachine", bigv.VirtualMachineName{VirtualMachine: "test-vm", Group: "test-group", Account: "test-account"}).Return(&vm, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ShowVM([]string{"test-vm.test-group.test-account"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestLockHWProfileCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}

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
	c := &mockBigVClient{}
	config := &mockConfig{}

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

func TestSetHWProfileCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}

	vmname := bigv.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account"}
	args := []string{"test-vm.test-group.test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return(args)

	// test no arguments, nothing should happen
	c.When("ParseVirtualMachineName", args[0]).Return(vmname)
	c.When("AuthWithToken", "test-token").Return(nil).Times(0)              // don't talk to BigV
	c.When("SetVirtualMachineHardwareProfile", vmname).Return(nil).Times(0) // don't do anything

	cmds := NewCommandSet(config, c)
	cmds.SetHWProfile(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test hardware profile only
	c.Reset()
	c.When("ParseVirtualMachineName", args[0]).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123").Return(nil).Times(1)

	args = []string{"test-vm.test-group.test-account", "virtio123"}

	cmds = NewCommandSet(config, c)
	cmds.SetHWProfile(args)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test --lock flag
	c.Reset()
	c.When("ParseVirtualMachineName", args[0]).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", true).Return(nil).Times(1)

	args_flag := []string{"test-vm.test-group.test-account", "virtio123", "--lock"}

	cmds = NewCommandSet(config, c)
	cmds.SetHWProfile(args_flag)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}

	// test --unlock flag
	c.Reset()
	c.When("ParseVirtualMachineName", args[0]).Return(vmname).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("SetVirtualMachineHardwareProfile", vmname, "virtio123", false).Return(nil).Times(1)

	args_flag = []string{"test-vm.test-group.test-account", "--unlock", "virtio123"}

	cmds = NewCommandSet(config, c)
	cmds.SetHWProfile(args_flag)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
