package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

//CreateVirtualMachine creates a virtual machine in the given group.
func (c *bytemarkClient) CreateVirtualMachine(group *GroupName, spec brain.VirtualMachineSpec) (vm *brain.VirtualMachine, err error) {
	err = c.validateGroupName(group)
	if err != nil {
		return nil, err
	}
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/vm_create", group.Account, group.Group)
	if err != nil {
		return nil, err
	}

	req := make(map[string]interface{})
	rvm := make(map[string]interface{})
	rvm["autoreboot_on"] = spec.VirtualMachine.Autoreboot
	if spec.VirtualMachine.CdromURL != "" {
		rvm["cdrom_url"] = spec.VirtualMachine.CdromURL
	}
	rvm["cores"] = spec.VirtualMachine.Cores
	rvm["memory"] = spec.VirtualMachine.Memory
	rvm["name"] = spec.VirtualMachine.Name
	if spec.VirtualMachine.HardwareProfile != "" {
		rvm["hardware_profile"] = spec.VirtualMachine.HardwareProfile
	}
	rvm["hardware_profile_locked"] = spec.VirtualMachine.HardwareProfileLocked
	if spec.VirtualMachine.ZoneName != "" {
		rvm["zone_name"] = spec.VirtualMachine.ZoneName
	}

	req["virtual_machine"] = rvm

	labelDiscs(spec.Discs)

	discs := make([]map[string]interface{}, 0, 4)

	for _, d := range spec.Discs {
		disc := make(map[string]interface{})
		label := d.Label
		disc["label"] = label
		disc["size"] = d.Size
		disc["storage_grade"] = d.StorageGrade

		discs = append(discs, disc)
	}

	req["discs"] = discs

	if spec.Reimage != nil {
		reimage := make(map[string]interface{})

		if spec.Reimage.Distribution != "" {
			reimage["distribution"] = spec.Reimage.Distribution
		}
		if spec.Reimage.RootPassword != "" {
			reimage["root_password"] = spec.Reimage.RootPassword
		}
		reimage["ssh_public_key"] = spec.Reimage.PublicKeys
		reimage["firstboot_script"] = spec.Reimage.FirstbootScript

		req["reimage"] = reimage
	}

	if spec.IPs != nil {
		ips := make(map[string]interface{})
		if spec.IPs.IPv4 != "" {
			ips["ipv4"] = spec.IPs.IPv4
		}
		if spec.IPs.IPv6 != "" {
			ips["ipv6"] = spec.IPs.IPv6
		}
		rvm["ips"] = ips
	}

	js, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	vm = new(brain.VirtualMachine)
	_, _, err = r.Run(bytes.NewBuffer(js), vm)
	return vm, err
}

// DeleteVirtualMachine deletes the named virtual machine.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) DeleteVirtualMachine(name *VirtualMachineName, purge bool) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	purgePart := ""
	if purge {
		purgePart = "?purge=true"
	}
	r, err := c.BuildRequest("DELETE", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s"+purgePart, name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.Run(nil, nil)
	return err
}

// GetVirtualMachine requests an overview of the named VM, regardless of its deletion status.
func (c *bytemarkClient) GetVirtualMachine(name *VirtualMachineName) (vm *brain.VirtualMachine, err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return
	}
	vm = new(brain.VirtualMachine)
	r, err := c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s?include_deleted=true&view=overview", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, vm)
	if err != nil {
		return
	}
	return
}

//MoveVirtualMachine moves the virtual machine to the given name, across groups if needed.
func (c *bytemarkClient) MoveVirtualMachine(oldName *VirtualMachineName, newName *VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(oldName)
	if err != nil {
		return
	}
	err = c.validateVirtualMachineName(newName)
	if err != nil {
		return
	}

	// create the change we want to see in the server
	change := brain.VirtualMachine{Name: newName.VirtualMachine}
	if newName.Group != "" || newName.Account != "" {
		// get group
		groupName := GroupName{Group: newName.Group, Account: newName.Account}
		group, err := c.GetGroup(&groupName)
		if err != nil {
			return err
		}
		change.GroupID = group.ID
	}

	// PUT the change
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", oldName.Account, oldName.Group, oldName.VirtualMachine)
	if err != nil {
		return err
	}

	js, err := json.Marshal(change)
	if err != nil {
		return err
	}
	_, _, err = r.Run(bytes.NewBuffer(js), nil)
	return err

}

// ReimageVirtualMachine reimages the named virtual machine. This will wipe everything on the first disk in the vm and install a new OS on top of it.
// Note that the machine in question must already be powered off. Once complete, according to the API docs, the vm will be powered on but its autoreboot_on will be false.
func (c *bytemarkClient) ReimageVirtualMachine(name *VirtualMachineName, image *brain.ImageInstall) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/reimage", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	js, err := json.Marshal(image)
	if err != nil {
		return err
	}
	_, _, err = r.Run(bytes.NewBuffer(js), nil)
	return err
}

// ResetVirtualMachine resets the named virtual machine. This is like pressing the reset
// button on a physical computer. This does not cause a new process to be started, so does not apply any pending hardware changes.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) ResetVirtualMachine(name *VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/signal", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.Run(bytes.NewBufferString(`{"signal":"reset"}`), nil)
	return err
}

// RestartVirtualMachine restarts the named virtual machine. This is
// returns nil on success or an error otherwise.
func (c *bytemarkClient) RestartVirtualMachine(name *VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.Run(bytes.NewBufferString(`{"autoreboot_on":true, "power_on": false}`), nil)
	return err
}

// StartVirtualMachine starts the named virtual machine.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) StartVirtualMachine(name *VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.Run(bytes.NewBufferString(`{"autoreboot_on":true, "power_on": true}`), nil)
	return err
}

// StopVirtualMachine starts the named virtual machine.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) StopVirtualMachine(name *VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.Run(bytes.NewBufferString(`{"autoreboot_on":false, "power_on": false}`), nil)
	return err
}

// ShutdownVirtualMachine sends an ACPI shutdown to the VM. This will cause a graceful shutdown of the machine
// returns nil on success or an error otherwise.
func (c *bytemarkClient) ShutdownVirtualMachine(name *VirtualMachineName, stayoff bool) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return
	}
	var r *Request
	if stayoff {
		r, err = c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
		if err != nil {
			return
		}

		_, _, err = r.Run(bytes.NewBufferString(`{"autoreboot_on":false}`), nil)
		if err != nil {
			return
		}
	}
	r, err = c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/signal", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return
	}

	_, _, err = r.Run(bytes.NewBufferString(`{"signal": "powerdown"}`), nil)
	return err
}

// UndeleteVirtualMachine changes the deleted flag on a VM back to false.
// Return nil on success, an error otherwise.
func (c *bytemarkClient) UndeleteVirtualMachine(name *VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.Run(bytes.NewBufferString(`{"deleted":false}`), nil)
	return err
}

// SetVirtualMachineHardwareProfile specifies the hardware profile on a VM. Optionally locks or unlocks h. profile
// Return nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineHardwareProfile(name *VirtualMachineName, profile string, locked ...bool) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}
	hwprofileLock := ""
	if len(locked) > 0 {
		hwprofileLock = `, "hardware_profile_locked": false`
		if locked[0] {
			hwprofileLock = `, "hardware_profile_locked": true`
		}
	}
	profileJSON := fmt.Sprintf(`{"hardware_profile": "%s"%s}`, profile, hwprofileLock)

	_, _, err = r.Run(bytes.NewBufferString(profileJSON), nil)
	return err
}

// SetVirtualMachineHardwareProfileLock locks or unlocks the hardware profile of a VM.
// Return nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineHardwareProfileLock(name *VirtualMachineName, locked bool) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	lockJSON := `{"hardware_profile_locked": false}`
	if locked {
		lockJSON = `{"hardware_profile_locked": true}`
	}

	_, _, err = r.Run(bytes.NewBufferString(lockJSON), nil)
	return err
}

// SetVirtualMachineMemory sets the RAM available to a virtual machine in megabytes
// Return nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineMemory(name *VirtualMachineName, memory int) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	memoryJSON := fmt.Sprintf(`{"memory": %d}`, memory)

	_, _, err = r.Run(bytes.NewBufferString(memoryJSON), nil)
	return err
}

// SetVirtualMachineCores sets the number of CPUs available to a virtual machine
// Return nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineCores(name *VirtualMachineName, cores int) (err error) {
	err = c.validateVirtualMachineName(name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	coresJSON := fmt.Sprintf(`{"cores": %d}`, cores)

	_, _, err = r.Run(bytes.NewBufferString(coresJSON), nil)
	return err
}
