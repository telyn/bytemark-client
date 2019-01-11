package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// CreateBackup creates a backup of the given disc, returning the backup if it was successful.
func (c *bytemarkClient) CreateBackup(path brain.DiscPather) (backup brain.Backup, err error) {
	path, err = c.checkDiscPather(path)
	if err != nil {
		return
	}

	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backups", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &backup)
	return
}

func (c *bytemarkClient) DeleteBackup(vm VirtualMachineName, discLabelOrID string, backupLabelOrID string) (err error) {
	err = c.checkVirtualMachinePather(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("DELETE", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backups/%s?purge=true", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID, backupLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, nil)
	return
}

func (c *bytemarkClient) GetBackups(vm VirtualMachineName, discLabelOrID string) (backups brain.Backups, err error) {
	err = c.checkVirtualMachinePather(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backups", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &backups)
	return
}

func (c *bytemarkClient) RestoreBackup(vm VirtualMachineName, discLabelOrID string, backupLabelOrID string) (backup brain.Backup, err error) {
	err = c.checkVirtualMachinePather(&vm)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backups/%s", vm.Account, vm.Group, vm.VirtualMachine, discLabelOrID, backupLabelOrID)
	if err != nil {
		return
	}
	restore := map[string]bool{
		"restore": true,
	}

	_, _, err = r.MarshalAndRun(restore, nil)
	return
}
