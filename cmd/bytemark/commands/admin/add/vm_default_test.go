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
		Name      string
		VMDefault brain.VirtualMachineDefault
		Args      []string
		Output    string
		ShouldErr bool
		// ResponseErr is the error returned by request.Run
		ResponseErr error
	}

	tests := []createTest{
		{
			Args: []string{
				"bytemark", "add", "vm", "default", "jeffrey",
			},
			ShouldErr: true,
		}, {
			Args: []string{
				"bytemark", "add", "vm", "default", "--image", "not-real-image", "jeffrey",
			},
			ShouldErr: true,
		}, {
			Name: "totally proper vm default",
			VMDefault: brain.VirtualMachineDefault{
				Name:   "test-vm-default",
				Public: true,
				ServerSettings: brain.VirtualMachineSpec{
					VirtualMachine: brain.VirtualMachine{
						Autoreboot:      true,
						CdromURL:        "https://example.com/example.iso",
						Cores:           1,
						HardwareProfile: "test-profile",
						Memory:          1024,
						PowerOn:         true,
						ZoneName:        "test-zone",
					},
					Discs: []brain.Disc{
						brain.Disc{
							Size:         25 * 1024,
							StorageGrade: "archive",
							BackupSchedules: brain.BackupSchedules{{
								Interval: 0,
								Capacity: 0,
							}},
						},
					},
					Reimage: &brain.ImageInstall{
						Distribution:    "test-image",
						FirstbootScript: "test-script",
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
				"test-vm-default", "true",
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
		t.Logf("TestCreateVMDefault %d", i)
		_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

		c.When("ReadDefinitions").Return(lib.Definitions{Distributions: []string{"test-image"}}, nil)

		request := mocks.Request{
			T:          t,
			StatusCode: 200,
			Err:        test.ResponseErr,
		}
		c.When("BuildRequest", "POST", lib.Endpoint(1), "/vm_defaults", []string(nil)).Return(&request)
		t.Logf("%#v", test)

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
	}
}
