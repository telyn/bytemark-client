package brain_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestStartVirtualMachineWithAppliance(t *testing.T) {
	tests := []struct {
		applianceName string
		vm            lib.VirtualMachineName
		expected      map[string]interface{}
		shouldErr     bool
	}{
		{
			applianceName: "rescue",
			vm: lib.VirtualMachineName{
				VirtualMachine: "test-vm",
				Group:          "test-group",
				Account:        "test-account",
			},
			expected: map[string]interface{}{
				"autoreboot_on": true,
				"power_on":      true,
				"appliance": map[string]interface{}{
					"name":      "rescue",
					"permanent": false},
			},
		},
	}
	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "PUT",
			Endpoint:      lib.BrainEndpoint,
			URL:           fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s", test.vm.Account, test.vm.Group, test.vm.VirtualMachine),
			AssertRequest: assert.BodyUnmarshalEqual(test.expected),
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := brainMethods.StartVirtualMachineWithAppliance(client, test.vm, test.applianceName)
			if test.shouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}
		})
	}
}
