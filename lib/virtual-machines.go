package lib

import (
	"encoding/json"
	"fmt"
)

//CreateVirtualMachine creates a virtual machine in the given group.
func (c *bytemarkClient) CreateVirtualMachine(group GroupName, spec VirtualMachineSpec) (vm *VirtualMachine, err error) {
	err = c.validateGroupName(&group)
	if err != nil {
		return nil, err
	}
	path := BuildURL("/accounts/%s/groups/%s/vm_create", group.Account, group.Group)

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

	vm = new(VirtualMachine)
	err = c.RequestAndUnmarshal(true, "POST", path, string(js), vm)
	return vm, err
}

// DeleteVirtualMachine deletes the named virtual machine.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) DeleteVirtualMachine(name VirtualMachineName, purge bool) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if purge {
		path += "?purge=true"
	}

	_, _, err = c.Request(true, "DELETE", path, "")
	return err
}

// GetVirtualMachine requests an overview of the named VM, regardless of its deletion status.
func (c *bytemarkClient) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return nil, err
	}
	vm = new(VirtualMachine)
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s?include_deleted=true&view=overview", name.Account, name.Group, name.VirtualMachine)

	err = c.RequestAndUnmarshal(true, "GET", path, "", vm)
	if err != nil {
		return nil, err
	}
	return vm, err
}

// ReimageVirtualMachine reimages the named virtual machine. This will wipe everything on the first disk in the vm and install a new OS on top of it.
// Note that the machine in question must already be powered off. Once complete, according to the API docs, the vm will be powered on but its autoreboot_on will be false.
func (c *bytemarkClient) ReimageVirtualMachine(name VirtualMachineName, image *ImageInstall) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/reimage", name.Account, name.Group, name.VirtualMachine)

	js, err := json.Marshal(image)
	if err != nil {
		return err
	}
	_, _, err = c.Request(true, "POST", path, string(js))
	return err
}

// ResetVirtualMachine resets the named virtual machine. This is like pressing the reset
// button on a physical computer. This does not cause a new process to be started, so does not apply any pending hardware changes.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) ResetVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/signal", name.Account, name.Group, name.VirtualMachine)

	_, _, err = c.Request(true, "POST", path, `{"signal":"reset"}`)
	return err
}

// RestartVirtualMachine restarts the named virtual machine. This is
// returns nil on success or an error otherwise.
func (c *bytemarkClient) RestartVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = c.Request(true, "PUT", path, `{"autoreboot_on":true, "power_on": false}`)
	return err
}

// StartVirtualMachine starts the named virtual machine.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) StartVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = c.Request(true, "PUT", path, `{"autoreboot_on":true, "power_on": true}`)
	return err
}

// StopVirtualMachine starts the named virtual machine.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) StopVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = c.Request(true, "PUT", path, `{"autoreboot_on":false, "power_on": false}`)
	return err
}

// ShutdownVirtualMachine sends an ACPI shutdown to the VM. This will cause a graceful shutdown of the machine
// returns nil on success or an error otherwise.
func (c *bytemarkClient) ShutdownVirtualMachine(name VirtualMachineName, stayoff bool) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	if stayoff {
		path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

		_, _, err = c.Request(true, "PUT", path, `{"autoreboot_on":false}`)
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/signal", name.Account, name.Group, name.VirtualMachine)

	_, _, err = c.Request(true, "POST", path, `{"signal": "powerdown"}`)
	return err
}

// UndeleteVirtualMachine changes the deleted flag on a VM back to false.
// Return nil on success, an error otherwise.
func (c *bytemarkClient) UndeleteVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = c.Request(true, "PUT", path, `{"deleted":false}`)
	return err
}

// SetVirtualMachineHardwareProfile specifies the hardware profile on a VM. Optionally locks or unlocks h. profile
// Return nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineHardwareProfile(name VirtualMachineName, profile string, locked ...bool) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	hwprofile_lock := ""
	if len(locked) > 0 {
		hwprofile_lock = `, "hardware_profile_locked": false`
		if locked[0] {
			hwprofile_lock = `, "hardware_profile_locked": true`
		}
	}
	what := fmt.Sprintf(`{"hardware_profile": "%s"%s}`, profile, hwprofile_lock)

	_, _, err = c.Request(true, "PUT", path, what)
	return err
}

// SetVirtualMachineHardwareProfileLock locks or unlocks the hardware profile of a VM.
// Return nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineHardwareProfileLock(name VirtualMachineName, locked bool) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	what := `{"hardware_profile_locked": false}`
	if locked {
		what = `{"hardware_profile_locked": true}`
	}

	_, _, err = c.Request(true, "PUT", path, what)
	return err
}

// SetVirtualMachineMemory sets the RAM available to a virtual machine in megabytes
// Return nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineMemory(name VirtualMachineName, memory int) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	what := fmt.Sprintf(`{"memory": %d}`, memory)

	_, _, err = c.Request(true, "PUT", path, what)
	return err
}

// SetVirtualMachineCores sets the number of CPUs available to a virtual machine
// Return nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineCores(name VirtualMachineName, cores int) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	what := fmt.Sprintf(`{"cores": %d}`, cores)

	_, _, err = c.Request(true, "PUT", path, what)
	return err
}
