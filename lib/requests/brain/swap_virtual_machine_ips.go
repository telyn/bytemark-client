package brain

import "github.com/BytemarkHosting/bytemark-client/lib"

// SwapVirtualMachineIPs swaps all primary IPs (v4 and v6) between two VMs. Both
// VMs must be powered off. If moveAdditional is provided, the additional IPs
// are also swapped between the machines.
// To swap only additional IPs, use multiple RerouteIP calls.. if we've written
// that already.
func SwapVirtualMachineIPs(client lib.Client, vmA lib.VirtualMachineName, vmB lib.VirtualMachineName, moveAdditional bool) error {
	a, err := client.GetVirtualMachine(vmA)
	if err != nil {
		return err
	}
	b, err := client.GetVirtualMachine(vmB)
	if err != nil {
		return err
	}
	r, err := client.BuildRequest("POST", lib.BrainEndpoint, "/ips/swap_virtual_machine_ips")
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"virtual_machine_1_id": a.ID,
		"virtual_machine_2_id": b.ID,
		"move_additional_ips":  moveAdditional,
	}

	_, _, err = r.MarshalAndRun(body, nil)
	return err
}
