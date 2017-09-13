package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"github.com/urfave/cli"
)

func TestCreateDiskCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("GetVirtualMachine").Return(&defVM)

	name := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "default", Account: "default-account"}
	c.When("GetVirtualMachine", &name).Return(&brain.VirtualMachine{Hostname: "test-server.default.default-account.endpoint"})

	disc := brain.Disc{Size: 35 * 1024, StorageGrade: "archive"}

	c.When("CreateDisc", &name, disc).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark create disc --force --disc archive:35 test-server", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateGroupCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("GetGroup").Return(&defGroup)

	group := lib.GroupName{
		Group:   "test-group",
		Account: "default-account",
	}
	c.When("CreateGroup", &group).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark create group test-group", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateServerHasCorrectFlags(t *testing.T) {
	is := is.New(t)
	seenCmd := false
	seenAuthKeys := false
	seenAuthKeysFile := false
	seenFirstbootScript := false
	seenFirstbootScriptFile := false
	seenImage := false
	seenRootPassword := false

	traverseAllCommands(commands, func(cmd cli.Command) {
		if cmd.FullName() == "create server" {
			seenCmd = true
			for _, f := range cmd.Flags {
				switch f.GetName() {
				case "authorized-keys":
					seenAuthKeys = true
				case "authorized-keys-file":
					seenAuthKeysFile = true
				case "firstboot-script":
					seenFirstbootScript = true
				case "firstboot-script-file":
					seenFirstbootScriptFile = true
				case "image":
					seenImage = true
				case "root-password":
					seenRootPassword = true
				}
			}
		}
	})
	is.True(seenCmd)
	is.True(seenAuthKeys)
	is.True(seenAuthKeysFile)
	is.True(seenFirstbootScript)
	is.True(seenFirstbootScriptFile)
	is.True(seenImage)
	is.True(seenRootPassword)

}

func TestCreateServer(t *testing.T) {
	type createTest struct {
		Spec                 brain.VirtualMachineSpec
		ConfigVirtualMachine lib.VirtualMachineName
		GroupName            lib.GroupName
		Args                 []string
		ShouldErr            bool
	}

	tests := []createTest{
		{
			Spec: brain.VirtualMachineSpec{
				Discs: []brain.Disc{
					brain.Disc{
						Size:         25 * 1024,
						StorageGrade: "sata",
					},
					brain.Disc{
						Size:         50 * 1024,
						StorageGrade: "archive",
					},
				},
				VirtualMachine: &brain.VirtualMachine{
					Name:                  "test-server",
					Autoreboot:            true,
					Cores:                 1,
					Memory:                1024,
					CdromURL:              "https://example.com/example.iso",
					HardwareProfile:       "test-profile",
					HardwareProfileLocked: true,
					ZoneName:              "test-zone",
				},
				Reimage: &brain.ImageInstall{
					Distribution:    "test-image",
					RootPassword:    "test-password",
					PublicKeys:      "test-pubkey",
					FirstbootScript: "test-script",
				},
				IPs: &brain.IPSpec{
					IPv4: "192.168.1.123",
					IPv6: "fe80::123",
				},
			},
			ConfigVirtualMachine: lib.VirtualMachineName{Group: "default"},
			GroupName:            lib.GroupName{Group: "default"},
			Args: []string{
				"bytemark", "create", "server",
				"--authorized-keys", "test-pubkey",
				"--firstboot-script", "test-script",
				"--cdrom", "https://example.com/example.iso",
				"--cores", "1",
				"--disc", "25",
				"--disc", "archive:50",
				"--force",
				"--hwprofile", "test-profile",
				"--hwprofile-locked",
				"--image", "test-image",
				"--ip", "192.168.1.123",
				"--ip", "fe80::123",
				"--memory", "1",
				"--root-password", "test-password",
				"--zone", "test-zone",
				"test-server",
			},
		},
		{
			ConfigVirtualMachine: defVM,
			Spec: brain.VirtualMachineSpec{
				VirtualMachine: &brain.VirtualMachine{
					Name:   "test-server",
					Cores:  1,
					Memory: 1024,
				},
				Discs: []brain.Disc{
					brain.Disc{
						Size:         25600,
						StorageGrade: "sata",
					},
				},
			},

			GroupName: lib.GroupName{
				Group:   "default",
				Account: "default-account",
			},
			Args: []string{
				"bytemark", "create", "server",
				"--cores", "1",
				"--force",
				"--memory", "1",
				"--no-image",
				"test-server",
			},
		},
		{
			ConfigVirtualMachine: defVM,
			GroupName: lib.GroupName{
				Group:   "default",
				Account: "default-account",
			},

			Spec: brain.VirtualMachineSpec{
				VirtualMachine: &brain.VirtualMachine{
					Name:   "test-server",
					Cores:  3,
					Memory: 6565,
				},
				Discs: []brain.Disc{{
					Size:         34 * 1024,
					StorageGrade: "archive",
				}},
			},
			Args: []string{
				"bytemark", "create", "server",
				"--force",
				"--no-image",
				"test-server", "3", "6565m", "archive:34",
			},
		},
	}

	var i int
	var test createTest
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestCreateVirtualMachine %d panicked.\r\n%v\r\n%v", i, r, string(debug.Stack()))
		}
	}()

	for i, test = range tests {
		t.Logf("TestCreateVirtualMachine %d", i)
		config, c := baseTestAuthSetup(t, false)
		config.When("GetVirtualMachine").Return(&test.ConfigVirtualMachine)

		vmname := lib.VirtualMachineName{
			VirtualMachine: test.Spec.VirtualMachine.Name,
			Group:          test.GroupName.Group,
			Account:        test.GroupName.Account,
		}

		discs := make([]*brain.Disc, len(test.Spec.Discs))
		for i, d := range test.Spec.Discs {
			discs[i] = &d
		}
		getvm := *test.Spec.VirtualMachine
		getvm.Discs = discs
		getvm.Hostname = "test-server.test-group.test-account.tld"

		c.When("CreateVirtualMachine", &test.GroupName, test.Spec).Return(test.Spec.VirtualMachine, nil).Times(1)
		c.When("GetVirtualMachine", &vmname).Return(&getvm, nil).Times(1)

		err := global.App.Run(test.Args)
		if err != nil {
			t.Error(err)
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}

func TestCreateBackup(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("GetVirtualMachine").Return(&defVM)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	c.When("CreateBackup", vmname, "test-disc").Return(brain.Backup{}, nil).Times(1)

	err := global.App.Run([]string{
		"bytemark", "create", "backup", "test-server", "test-disc",
	})
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestCreateVLANGroup(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetGroup").Return(&defGroup).Times(1)

	group := lib.GroupName{
		Group:   "test-group",
		Account: "test-account",
	}
	c.When("AdminCreateGroup", group, 0).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark create vlan_group test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateVLANGroupWithVLANNum(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetGroup").Return(&defGroup).Times(1)

	group := lib.GroupName{
		Group:   "test-group",
		Account: "test-account",
	}
	c.When("AdminCreateGroup", group, 19).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark create vlan_group test-group.test-account 19", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateVLANGroupError(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, true)

	config.When("GetGroup").Return(&defGroup).Times(1)

	group := lib.GroupName{
		Group:   "test-group",
		Account: "test-account",
	}
	c.When("AdminCreateGroup", group, 0).Return(fmt.Errorf("Group name already used")).Times(1)

	err := global.App.Run(strings.Split("bytemark create vlan_group test-group.test-account", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateIPRange(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CreateIPRange", "192.168.3.0/28", 14).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark create ip_range 192.168.3.0/28 14", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateIPRangeError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CreateIPRange", "192.168.3.0/28", 18).Return(fmt.Errorf("Error creating IP range")).Times(1)

	err := global.App.Run(strings.Split("bytemark create ip_range 192.168.3.0/28 18", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CreateUser", "uname", "cluster_su").Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark create user uname cluster_su", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateUserError(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	c.When("CreateUser", "uname", "cluster_su").Return(fmt.Errorf("Error creating user")).Times(1)

	err := global.App.Run(strings.Split("bytemark create user uname cluster_su", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
