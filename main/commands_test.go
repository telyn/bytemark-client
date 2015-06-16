package main

import (
	bigv "bigv.io/client/lib"
	"testing"
	//"github.com/cheekybits/is"
)

func getFixtureVM() bigv.VirtualMachine {
	return bigv.VirtualMachine{
		Name:    "test-vm",
		GroupID: 1,
	}
}

func TestCommandConfig(t *testing.T) {
	config := &mockConfig{}

	config.When("GetV", "user").Return(ConfigVar{"user", "old-test-user", "config"})
	config.When("Get", "user").Return("old-test-user")
	config.When("Get", "silent").Return("true")

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
	config.When("GetBool", "force").Return(true)
	config.When("GetBool", "silent").Return(true)
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

func TestShowVMCommand(t *testing.T) {
	c := &mockBigVClient{}
	config := &mockConfig{}

	config.When("Get", "token").Return("test-token")
	config.When("Get", "silent").Return("true")
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
