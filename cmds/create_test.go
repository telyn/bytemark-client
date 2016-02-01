package cmds

import (
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestCreateDiskCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}
	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{"", "", ""})

	config.When("ImportFlags").Return([]string{"test-vm", "archive:35"})
	name := lib.VirtualMachineName{VirtualMachine: "test-vm"}
	c.When("ParseVirtualMachineName", "test-vm", []lib.VirtualMachineName{{}}).Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	c.When("GetVirtualMachine", name).Return(&lib.VirtualMachine{Hostname: "test-vm.default.test-user.endpoint"})

	disc := lib.Disc{Size: 35 * 1024, StorageGrade: "archive"}

	c.When("CreateDisc", name, disc).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.CreateDiscs([]string{"test-vm", "archive:35"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

// TODO(telyn): TestCreateGroupCommand

func TestCreateVMCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-vm"})
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{"", "", ""})

	c.When("ParseVirtualMachineName", "test-vm", []lib.VirtualMachineName{{}}).Return(lib.VirtualMachineName{VirtualMachine: "test-vm"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	vm := lib.VirtualMachineSpec{
		Discs: []lib.Disc{
			lib.Disc{
				Size:         25 * 1024,
				StorageGrade: "sata",
			},
			lib.Disc{
				Size:         50 * 1024,
				StorageGrade: "archive",
			},
		},
		VirtualMachine: &lib.VirtualMachine{
			Name:                  "test-vm",
			Autoreboot:            true,
			Cores:                 1,
			Memory:                1024,
			CdromURL:              "https://example.com/example.iso",
			HardwareProfile:       "test-profile",
			HardwareProfileLocked: true,
			ZoneName:              "test-zone",
		},
		Reimage: &lib.ImageInstall{
			Distribution: "test-image",
			RootPassword: "test-password",
		},
		IPs: &lib.IPSpec{
			IPv4: "192.168.1.123",
			IPv6: "fe80::123",
		},
	}

	group := lib.GroupName{
		Group:   "",
		Account: "",
	}

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "",
		Account:        "",
	}

	c.When("CreateVirtualMachine", group, vm).Return(vm, nil).Times(1)
	c.When("GetVirtualMachine", vmname).Return(vm.VirtualMachine, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.CreateVM([]string{
		"--cdrom", "https://example.com/example.iso",
		"--cores", "1",
		"--disc", "25",
		"--disc", "archive:50",
		"--hwprofile", "test-profile",
		"--hwprofile-locked",
		"--image", "test-image",
		"--ip", "192.168.1.123",
		"--ip", "fe80::123",
		"--memory", "1",
		"--root-password", "test-password",
		"--zone", "test-zone",
		"test-vm",
	})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateVMNoImagesNoDiscs(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-vm"})
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{"", "", ""})

	c.When("ParseVirtualMachineName", "test-vm", []lib.VirtualMachineName{{}}).Return(lib.VirtualMachineName{VirtualMachine: "test-vm"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	vm := lib.VirtualMachineSpec{
		VirtualMachine: &lib.VirtualMachine{
			Name:   "test-vm",
			Cores:  1,
			Memory: 1024,
		},
	}

	group := lib.GroupName{
		Group:   "",
		Account: "",
	}

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "",
		Account:        "",
	}

	c.When("CreateVirtualMachine", group, vm).Return(vm, nil).Times(1)
	c.When("GetVirtualMachine", vmname).Return(vm.VirtualMachine, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.CreateVM([]string{
		"--cores", "1",
		"--no-discs",
		"--memory", "1",
		"test-vm",
	})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateVM(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-vm", "3", "6565m", "archive:34"})
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{"", "", ""})

	c.When("ParseVirtualMachineName", "test-vm", []lib.VirtualMachineName{{}}).Return(lib.VirtualMachineName{VirtualMachine: "test-vm"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	vm := lib.VirtualMachineSpec{
		VirtualMachine: &lib.VirtualMachine{
			Name:   "test-vm",
			Cores:  3,
			Memory: 6565,
		},
		Discs: []lib.Disc{{
			Size:         34 * 1024,
			StorageGrade: "archive",
		},
		},
	}

	group := lib.GroupName{
		Group:   "",
		Account: "",
	}

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "",
		Account:        "",
	}

	c.When("CreateVirtualMachine", group, vm).Return(vm, nil).Times(1)
	c.When("GetVirtualMachine", vmname).Return(vm.VirtualMachine, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.CreateVM([]string{
		"--no-image",
		"test-vm", "3", "6565m", "archive:34",
	})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
