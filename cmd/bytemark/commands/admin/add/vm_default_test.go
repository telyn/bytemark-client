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
		Public    bool
		Spec      brain.VMDefaultSpec
		Args      []string
		Output    string
		ShouldErr bool
	}

	tests := []createTest{
		{
			Name:   "vmdefault",
			Public: true,
			Spec: brain.VMDefaultSpec{
				VMDefault: brain.VMDefault{
					CdromURL:        "https://example.com/example.iso",
					Cores:           1,
					Memory:          1024,
					Name:            "test-vm",
					HardwareProfile: "test-profile",
					ZoneName:        "test-zone",
					Discs:           nil,
					ID:              0,
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
			Args: []string{
				"bytemark", "add", "vm default", "vmdefault", "true",
				"--cdrom", "https://example.com/example.iso",
				"--cores", "1",
				"--memory", "1",
				"--vm-name", "vm",
				"--hwprofile", "test-profile",
				"--backup", "never",
				"--zone", "test-zone",
				"--disc", "archive:50",
				"--image", "test-image",
				"--firstboot-script", "test-script",
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

		c.When("BuildRequest", "POST", lib.Endpoint(1), "/vm_defaults", []string(nil)).Return(&mocks.Request{
			T:              t,
			StatusCode:     200,
			ResponseObject: nil,
		})

		c.When("CreateVMDefault", test.Name, test.Public, test.Spec).Return(nil).Times(1)

		c.When("ReadDefinitions").Return(lib.Definitions{Distributions: []string{"test-image"}}, nil)

		err := app.Run(test.Args)
		if err != nil {
			t.Error(err)
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}
