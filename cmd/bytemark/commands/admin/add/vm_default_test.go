package add_test

import (
	"runtime/debug"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/mocks"

	"github.com/BytemarkHosting/bytemark-client/lib"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func TestCreateVMDefaultCommand(t *testing.T) {
	type createTest struct {
		Name string
		// account that will be requested to find out account id
		// (use VMDefault.AccountID to specify account id)
		Account string
		// VMDefault to expect
		VMDefault brain.VirtualMachineDefault
		// Args to provide to app.Run
		Args []string
		// Output to expect... one day we'll write output tests
		Output string
		// ShouldErr should be set to true if you expect app.Run to error
		ShouldErr bool
		// ResponseErr is the error to return from request.Run (to simulate 403/404/network error/etc)
		ResponseErr error
	}

	tests := []createTest{
		{
			Name:    "no params",
			Account: "bytemark",
			Args: []string{
				"bytemark", "add", "vm", "default", "jeffrey",
			},
			ShouldErr: true,
		}, {
			Name:    "nonexistent image",
			Account: "bytemark",
			Args: []string{
				"bytemark", "add", "vm", "default", "--image", "not-real-image", "jeffrey",
			},
			ShouldErr: true,
		}, {
			Name:    "public on bytemark account",
			Account: "bytemark",
			VMDefault: brain.VirtualMachineDefault{
				AccountID: 142,
				Name:      "test-vm-default",
				Public:    true,
				ServerSettings: brain.VirtualMachineSpec{
					VirtualMachine: brain.VirtualMachine{
						Autoreboot:      true,
						CdromURL:        "https://example.com/example.iso",
						Cores:           1,
						HardwareProfile: "test-profile",
						Memory:          1024,
						ZoneName:        "test-zone",
					},
					Discs: []brain.Disc{
						brain.Disc{
							Size:            50 * 1024,
							StorageGrade:    "archive",
							BackupSchedules: nil,
						},
					},
					Reimage: &brain.ImageInstall{
						Distribution:    "test-image",
						FirstbootScript: "test-script",
						RootPassword:    "test-password",
					},
				},
			},
			Args: []string{
				"bytemark", "add", "vm", "default",
				"--cdrom", "https://example.com/example.iso",
				"--cores", "1",
				"--memory", "1",
				"--hwprofile", "test-profile",
				"--backup", "never",
				"--zone", "test-zone",
				"--disc", "archive:50",
				"--image", "test-image",
				"--firstboot-script", "test-script",
				"--root-password", "test-password",
				"--public",
				"test-vm-default", "true",
			},
		}, {
			Name:    "private on bytemark account",
			Account: "bytemark",
			VMDefault: brain.VirtualMachineDefault{
				AccountID: 142,
				Name:      "test-vm-default",
				Public:    false,
				ServerSettings: brain.VirtualMachineSpec{
					VirtualMachine: brain.VirtualMachine{
						Autoreboot:      true,
						CdromURL:        "https://example.com/example.iso",
						Cores:           1,
						HardwareProfile: "test-profile",
						Memory:          1024,
						ZoneName:        "test-zone",
					},
					Discs: []brain.Disc{
						brain.Disc{
							Size:            50 * 1024,
							StorageGrade:    "archive",
							BackupSchedules: nil,
						},
					},
					Reimage: &brain.ImageInstall{
						Distribution:    "test-image",
						FirstbootScript: "test-script",
						RootPassword:    "test-password",
					},
				},
			},
			Args: []string{
				"bytemark", "add", "vm", "default",
				"--cdrom", "https://example.com/example.iso",
				"--cores", "1",
				"--memory", "1",
				"--hwprofile", "test-profile",
				"--backup", "never",
				"--zone", "test-zone",
				"--disc", "archive:50",
				"--image", "test-image",
				"--firstboot-script", "test-script",
				"--root-password", "test-password",
				"test-vm-default", "true",
			},
		}, {
			Name:    "public on tomatoes account",
			Account: "tomatoes",
			VMDefault: brain.VirtualMachineDefault{
				AccountID: 26580,
				Name:      "test-vm-default",
				Public:    true,
				ServerSettings: brain.VirtualMachineSpec{
					VirtualMachine: brain.VirtualMachine{
						Autoreboot:      true,
						CdromURL:        "https://example.com/example.iso",
						Cores:           1,
						HardwareProfile: "test-profile",
						Memory:          1024,
						ZoneName:        "test-zone",
					},
					Discs: []brain.Disc{
						brain.Disc{
							Size:            50 * 1024,
							StorageGrade:    "archive",
							BackupSchedules: nil,
						},
					},
					Reimage: &brain.ImageInstall{
						Distribution:    "test-image",
						FirstbootScript: "test-script",
						RootPassword:    "test-password",
					},
				},
			},
			Args: []string{
				"bytemark", "add", "vm", "default",
				"--cdrom", "https://example.com/example.iso",
				"--cores", "1",
				"--memory", "1",
				"--hwprofile", "test-profile",
				"--backup", "never",
				"--zone", "test-zone",
				"--disc", "archive:50",
				"--image", "test-image",
				"--firstboot-script", "test-script",
				"--root-password", "test-password",
				"--account", "tomatoes",
				"--public",
				"test-vm-default",
			},
		},
	}

	var i int
	var test createTest
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestCreateVMDefault %d panicked.\r\n%v\r\n%v", i, r, string(debug.Stack()))
		}
	}()

	for i, test = range tests {
		t.Run(test.Name, func(t *testing.T) {
			config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

			config.When("GetIgnoreErr", "account").Return("jeff")

			if test.Account != "" {
				c.When("GetAccount", test.Account).Return(lib.Account{BrainID: test.VMDefault.AccountID})
			}

			c.When("ReadDefinitions").Return(lib.Definitions{Distributions: []string{"test-image"}}, nil)

			request := mocks.Request{
				T:          t,
				StatusCode: 200,
				Err:        test.ResponseErr,
			}
			c.When("BuildRequest", "POST", lib.Endpoint(1), "/vm_defaults", []string(nil)).Return(&request)

			err := app.Run(test.Args)
			if !test.ShouldErr && err != nil {
				t.Errorf("Unexpected error: %s", err)
			} else if test.ShouldErr && err == nil {
				t.Error("Expected error but didn't get one")
			}

			if !test.ShouldErr {
				request.AssertRequestObjectEqual(test.VMDefault)
			}
			if ok, err := c.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}
