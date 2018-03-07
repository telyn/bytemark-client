package brain

import (
	"bytes"
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
)

// StartVirtualMachineWithAppliance starts the named virtual machine using the named appliance.
// returns nil on success or an error otherwise.
func StartVirtualMachineWithAppliance(client lib.Client, vmName lib.VirtualMachineName, applianceName string) (err error) {
	err = client.EnsureVirtualMachineName(&vmName)
	if err != nil {
		return err
	}
	r, err := client.BuildRequest("PUT", lib.BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", vmName.Account, vmName.Group, vmName.VirtualMachine)
	if err != nil {
		return err
	}

	body := fmt.Sprintf(`{"autoreboot_on":true, "power_on": true, "appliance":{"name":"%s", "permanent": false}}`, applianceName)

	_, _, err = r.Run(bytes.NewBufferString(body), nil)
	return err
}
