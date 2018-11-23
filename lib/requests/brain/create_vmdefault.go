package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// CreateVMDefault creates a new VM Default with the specified parameters,
// returning the newly created VM Default on success or an error otherwise.
func CreateVMDefault(client lib.Client, spec brain.VirtualMachineDefault) (created brain.VirtualMachineDefault, err error) {

	req, err := client.BuildRequest("POST", lib.BrainEndpoint, "/vm_defaults")
	if err != nil {
		return
	}

	_, _, err = req.MarshalAndRun(spec, &created)
	return
}
