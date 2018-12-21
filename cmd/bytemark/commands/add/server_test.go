package add_test

import (
	"runtime/debug"
	"testing"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/add"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
)

func TestCreateServerHasCorrectFlags(t *testing.T) {
	seenCmd := false
	seenAuthKeys := false
	seenAuthKeysFile := false
	seenFirstbootScript := false
	seenFirstbootScriptFile := false
	seenImage := false
	seenRootPassword := false

	testutil.TraverseAllCommands(add.Commands, func(cmd cli.Command) {
		if cmd.Name == "server" {
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
	if !seenCmd {
		t.Error("Didn't see add server command")
	}
	if !seenAuthKeys {
		t.Error("Didn't see authorized-keys")
	}
	if !seenAuthKeysFile {
		t.Error("Didn't see authorised-keys-file")
	}
	if !seenFirstbootScript {
		t.Error("Didn't see firstboot-script")
	}
	if !seenFirstbootScriptFile {
		t.Error("Didn't see firstboot-script-file")
	}
	if !seenImage {
		t.Error("Didn't see image")
	}
	if !seenRootPassword {
		t.Error("Didn't see root-password")
	}

}

func TestCreateServer(t *testing.T) {
	type createTest struct {
		Spec                 brain.VirtualMachineSpec
		ConfigVirtualMachine lib.VirtualMachineName
		GroupName            lib.GroupName
		Args                 []string
		Output               string
		ShouldErr            bool
	}

	tomorrow := time.Now().Add(24 * time.Hour)
	y, m, d := tomorrow.Date()
	midnightTonight := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	defaultStartDate := midnightTonight.Format("2006-01-02 15:04:05 MST")

	tests := []createTest{
		{
			Spec: brain.VirtualMachineSpec{
				Discs: []brain.Disc{
					brain.Disc{
						Size:         25 * 1024,
						StorageGrade: "sata",
						BackupSchedules: brain.BackupSchedules{{
							StartDate: defaultStartDate,
							Interval:  7 * 86400,
							Capacity:  1,
						}},
					},
					brain.Disc{
						Size:         50 * 1024,
						StorageGrade: "archive",
					},
				},
				VirtualMachine: brain.VirtualMachine{
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
				"bytemark", "add", "server",
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
		}, {
			ConfigVirtualMachine: testutil.DefVM,
			Spec: brain.VirtualMachineSpec{
				VirtualMachine: brain.VirtualMachine{
					Name:   "test-server",
					Cores:  1,
					Memory: 1024,
				},
				Discs: []brain.Disc{
					brain.Disc{
						Size:         25600,
						StorageGrade: "sata",
						BackupSchedules: brain.BackupSchedules{{
							StartDate: defaultStartDate,
							Interval:  7 * 86400,
							Capacity:  1,
						}},
					},
				},
			},

			GroupName: lib.GroupName{
				Group:   "default",
				Account: "default-account",
			},
			Args: []string{
				"bytemark", "add", "server",
				"--cores", "1",
				"--force",
				"--memory", "1",
				"--no-image",
				"test-server",
			},
		}, {
			ConfigVirtualMachine: testutil.DefVM,
			GroupName: lib.GroupName{
				Group:   "default",
				Account: "default-account",
			},

			Spec: brain.VirtualMachineSpec{
				VirtualMachine: brain.VirtualMachine{
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
				"bytemark", "add", "server",
				"--force",
				"--no-image",
				"--backup", "never",
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
		config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
		config.When("GetVirtualMachine").Return(test.ConfigVirtualMachine)

		vmname := lib.VirtualMachineName{
			VirtualMachine: test.Spec.VirtualMachine.Name,
			Group:          test.GroupName.Group,
			Account:        test.GroupName.Account,
		}

		postGpName := test.GroupName
		_ = c.EnsureGroupName(&postGpName)

		getvm := test.Spec.VirtualMachine
		getvm.Discs = test.Spec.Discs
		getvm.Hostname = "test-server.test-group.test-account.tld"

		c.When("CreateVirtualMachine", postGpName, test.Spec).Return(test.Spec.VirtualMachine, nil).Times(1)
		c.When("GetVirtualMachine", vmname).Return(getvm, nil).Times(1)

		err := app.Run(test.Args)
		if err != nil {
			t.Error(err)
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}
