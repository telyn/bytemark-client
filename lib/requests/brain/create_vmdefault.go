package brain

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// CreateVMDefault creates a new VM Default with the specified parameters,
// returning the newly created VM Default on success or an error otherwise.
func CreateVMDefault(client lib.Client, name string, public bool, serverSettings brain.VMDefaultSpec) (err error) {
	if name == "" {
		return fmt.Errorf("VMDefault must have a non-blank name")
	}
	if serverSettings.VMDefault.Name == "" {
		return fmt.Errorf("VM must have a non-blank name")
	}

	req, err := client.BuildRequest("POST", lib.BrainEndpoint, "/vm_defaults")
	if err != nil {
		return
	}

	obj := map[string]interface{}{
		"name":            name,
		"public":          public,
		"server_settings": serverSettings,
	}

	_, _, err = req.MarshalAndRun(obj, nil)
	return
}
