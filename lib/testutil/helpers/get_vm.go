package helpers

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
)

func GetVM(vmName lib.VirtualMachineName, vm brain.VirtualMachine) testutil.RequestTestSpec {
	return testutil.RequestTestSpec{
		Endpoint: lib.BrainEndpoint,
		Method:   "GET",
		URL:      fmt.Sprintf("/accounts/%s/groups/%s/virtual_machines/%s", vmName.Account, vmName.Group, vmName.VirtualMachine),
		Response: vm,
	}
}
