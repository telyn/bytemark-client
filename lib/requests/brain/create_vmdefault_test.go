package brain_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestCreateVMDefault(t *testing.T) {
	tests := []struct {
		Input     brain.VirtualMachineDefault
		Expect    brain.VirtualMachineDefault
		ExpectErr bool
	}{
		{
			brain.VirtualMachineDefault{
				Name:   "vm-default",
				Public: true,
				ServerSettings: brain.VirtualMachineSpec{
					VirtualMachine: brain.VirtualMachine{
						CdromURL:        "test-url",
						Cores:           1,
						Memory:          1,
						HardwareProfile: "test-profile",
						ZoneName:        "test-zone",
						Discs:           nil,
					},
					Discs: brain.Discs{
						brain.Disc{
							StorageGrade: "sata",
							Size:         1,
							BackupSchedules: brain.BackupSchedules{{
								Interval: 604800,
								Capacity: 1,
							}},
						},
					},
					Reimage: &brain.ImageInstall{
						Distribution:    "test-image",
						FirstbootScript: "test-script",
					},
				},
			},
			brain.VirtualMachineDefault{
				Name:   "vm-default",
				Public: true,
				ServerSettings: brain.VirtualMachineSpec{
					VirtualMachine: brain.VirtualMachine{
						CdromURL:        "test-url",
						Cores:           1,
						Memory:          1,
						HardwareProfile: "test-profile",
						ZoneName:        "test-zone",
						Discs:           nil,
					},
					Discs: brain.Discs{
						brain.Disc{
							StorageGrade: "sata",
							Size:         1,
							BackupSchedules: brain.BackupSchedules{{
								Interval: 604800,
								Capacity: 1,
							}},
						},
					},
					Reimage: &brain.ImageInstall{
						Distribution:    "test-image",
						FirstbootScript: "test-script",
					},
				},
			},
			false,
		},
	}
	for i, test := range tests {
		testName := testutil.Name(i)
		spec := brain.VirtualMachineDefault{}
		rts := testutil.RequestTestSpec{
			Method:   "POST",
			Endpoint: lib.BrainEndpoint,
			URL:      "/vm_defaults",
			AssertRequest: assert.BodyUnmarshal(&spec, func(_ *testing.T, _ string) {
				assert.Equal(t, testName, test.Expect, spec)
			}),
			Response: test.Expect.ServerSettings,
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := brainRequests.CreateVMDefault(client, test.Input)
			if err != nil && !test.ExpectErr {
				t.Fatal(err)
			}
		})
	}
}
