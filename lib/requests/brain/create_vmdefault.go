package brain

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func CreateVMDefault(client lib.Client, name string, public bool, serverSettings brain.VMDefaultSpec) error {
	if name == "" {
		return fmt.Errorf("VMDefault must have a non-blank name")
	}
	if serverSettings.VMDefault.Name == "" {
		return fmt.Errorf("VM must have a non-blank name")
	}

	r, err := client.BuildRequest("POST", lib.BrainEndpoint, "/vm_defaults")
	if err != nil {
		return err
	}

	obj := map[string]interface{}{
		"name":            name,
		"public":          public,
		"server_settings": serverSettings,
	}

	_, _, err = r.MarshalAndRun(obj, nil)
	return err
}
