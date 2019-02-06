package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// CreateBackup creates a backup of the given disc, returning the backup if it was successful.
func (c *bytemarkClient) CreateBackup(vm pathers.VirtualMachineName, discLabelOrID string) (backup brain.Backup, err error) {
	err = c.EnsureVirtualMachineName(&vm)
	if err != nil {
		return
	}

	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backups", string(vm.Account), vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &backup)
	return
}

func (c *bytemarkClient) DeleteBackup(vm pathers.VirtualMachineName, discLabelOrID string, backupLabelOrID string) (err error) {
	err = c.EnsureVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("DELETE", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backups/%s?purge=true", string(vm.Account), vm.Group, vm.VirtualMachine, discLabelOrID, backupLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) GetBackups(vm pathers.VirtualMachineName, discLabelOrID string) (backups brain.Backups, err error) {
	err = c.EnsureVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backups", string(vm.Account), vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &backups)
	return
}

func (c *bytemarkClient) RestoreBackup(vm pathers.VirtualMachineName, discLabelOrID string, backupLabelOrID string) (backup brain.Backup, err error) {
	err = c.EnsureVirtualMachineName(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backups/%s", string(vm.Account), vm.Group, vm.VirtualMachine, discLabelOrID, backupLabelOrID)
	if err != nil {
		return
	}
	restore := map[string]bool{
		"restore": true,
	}

	_, _, err = r.MarshalAndRun(restore, nil)
	return
}
