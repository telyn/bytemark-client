package brain_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/helpers"
)

func TestSwapVirtualMachineIPs(t *testing.T) {
	a := pathers.VirtualMachineName{
		VirtualMachine: "swappo",
		GroupName: pathers.GroupName{
			Group:   "test-group",
			Account: "bytemark",
		},
	}
	b := pathers.VirtualMachineName{
		VirtualMachine: "swappy",
		GroupName: pathers.GroupName{
			Group:   "test-group",
			Account: "bytemark",
		},
	}
	vmA := brain.VirtualMachine{
		ID: 1903,
	}
	vmB := brain.VirtualMachine{
		ID: 1124,
	}
	tests := []struct {
		name string

		nameA          pathers.VirtualMachineName
		nameB          pathers.VirtualMachineName
		vmA            brain.VirtualMachine
		vmB            brain.VirtualMachine
		moveAdditional bool

		expectedRequest map[string]interface{}
		shouldErr       bool
	}{
		{
			name:           "simple success",
			nameA:          a,
			nameB:          b,
			vmA:            vmA,
			vmB:            vmB,
			moveAdditional: false,
			expectedRequest: map[string]interface{}{
				"virtual_machine_1_id": float64(vmA.ID),
				"virtual_machine_2_id": float64(vmB.ID),
				"move_additional_ips":  false,
			},
		},
		{
			name:           "move additional",
			nameA:          a,
			nameB:          b,
			vmA:            vmA,
			vmB:            vmB,
			moveAdditional: true,
			expectedRequest: map[string]interface{}{
				"virtual_machine_1_id": float64(vmA.ID),
				"virtual_machine_2_id": float64(vmB.ID),
				"move_additional_ips":  true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mrts := testutil.MultiRequestTestSpec{
				Auth: true,
				Specs: []testutil.RequestTestSpec{
					helpers.GetVM(tc.nameA, tc.vmA),
					helpers.GetVM(tc.nameB, tc.vmB),
					testutil.RequestTestSpec{
						Endpoint:      lib.BrainEndpoint,
						Method:        "POST",
						URL:           "/ips/swap_virtual_machine_ips",
						AssertRequest: assert.BodyUnmarshalEqual(tc.expectedRequest),
					},
				},
			}
			mrts.Run(t, func(c lib.Client) {
				err := brainRequests.SwapVirtualMachineIPs(c, tc.nameA, tc.nameB, tc.moveAdditional)
				if err != nil && !tc.shouldErr {
					t.Errorf("No error expected, but %v received", err)
				} else if err == nil && tc.shouldErr {
					t.Error("Expected an error but SwapPrimaryIPs returned no error")
				}
			})
		})
	}
}
