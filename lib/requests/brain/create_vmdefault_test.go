package brain_test

import (
	"reflect"
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
			Response: brain.VirtualMachineDefault{
				Name: "jeff",
			},
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			vmd, err := brainRequests.CreateVMDefault(client, test.Input)
			if !reflect.DeepEqual(brain.VirtualMachineDefault{
				Name: "jeff",
			}, vmd) {
				t.Errorf("response wasn't jeff :-( got: %#v", vmd)
			}
			if err != nil && !test.ExpectErr {
				t.Fatal(err)
			}
		})
	}
}
