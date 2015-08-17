package lib

import (
	"encoding/json"
	"fmt"
)

var i2b = [...]string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

//CreateVirtualMachine creates a virtual machine in the given group.
func (bigv *bigvClient) CreateVirtualMachine(group GroupName, spec VirtualMachineSpec) (vm *VirtualMachine, err error) {
	err = bigv.validateGroupName(&group)
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

	discs := make([]map[string]interface{}, 0, 4)

	for i, d := range spec.Discs {
		if i > 7 {
			return nil, TooManyDiscsOnTheDancefloorError{}
		}
		disc := make(map[string]interface{})
		label := d.Label
		if label == "" {
			label = "vd" + i2b[i]
		}
		disc["label"] = label
		disc["size"] = d.Size
		disc["storage_grade"] = d.StorageGrade

		discs = append(discs, disc)
	}

	req["discs"] = discs

	reimage := make(map[string]interface{})

	if spec.Reimage.Distribution != "" {
		reimage["distribution"] = spec.Reimage.Distribution
	}
	if spec.Reimage.RootPassword != "" {
		reimage["root_password"] = spec.Reimage.RootPassword
	}
	reimage["ssh_public_key"] = spec.Reimage.PublicKeys

	req["reimage"] = reimage

	js, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	vm = new(VirtualMachine)
	err = bigv.RequestAndUnmarshal(true, "POST", path, string(js), vm)
	return vm, err
}

// DeleteVirtualMachine deletes the named virtual machine.
// returns nil on success or an error otherwise.
func (bigv *bigvClient) DeleteVirtualMachine(name VirtualMachineName, purge bool) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if purge {
		path += "?purge=true"
	}

	_, _, err = bigv.Request(true, "DELETE", path, "")
	return err
}

// GetVirtualMachine requests an overview of the named VM, regardless of its deletion status.
func (bigv *bigvClient) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return nil, err
	}
	vm = new(VirtualMachine)
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s?include_deleted=true&view=overview", name.Account, name.Group, name.VirtualMachine)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", vm)
	if err != nil {
		return nil, err
	}
	return vm, err
}

// ResetVirtualMachine resets the named virtual machine. This is like pressing the reset
// button on a physical computer. This does not cause a new process to be started, so does not apply any pending hardware changes.
// returns nil on success or an error otherwise.
func (bigv *bigvClient) ResetVirtualMachine(name VirtualMachineName) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/signal", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "POST", path, `{"signal":"reset"}`)
	return err
}

// RestartVirtualMachine restarts the named virtual machine. This is
// returns nil on success or an error otherwise.
func (bigv *bigvClient) RestartVirtualMachine(name VirtualMachineName) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"autoreboot_on":true, "power_on": false}`)
	return err
}

// StartVirtualMachine starts the named virtual machine.
// returns nil on success or an error otherwise.
func (bigv *bigvClient) StartVirtualMachine(name VirtualMachineName) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"autoreboot_on":true, "power_on": true}`)
	return err
}

// StopVirtualMachine starts the named virtual machine.
// returns nil on success or an error otherwise.
func (bigv *bigvClient) StopVirtualMachine(name VirtualMachineName) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"autoreboot_on":false, "power_on": false}`)
	return err
}

// ShutdownVirtualMachine sends an ACPI shutdown to the VM. This will cause a graceful shutdown of the machine
// returns nil on success or an error otherwise.
func (bigv *bigvClient) ShutdownVirtualMachine(name VirtualMachineName, stayoff bool) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	if stayoff {
		path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

		_, _, err = bigv.Request(true, "PUT", path, `{"autoreboot_on":false}`)
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s/signal", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"signal": "powerdown"}`)
	return err
}

// UndeleteVirtualMachine changes the deleted flag on a VM back to false.
// Return nil on success, an error otherwise.
func (bigv *bigvClient) UndeleteVirtualMachine(name VirtualMachineName) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"deleted":false}`)
	return err
}

// SetVirtualMachineHardwareProfile specifies the hardware profile on a VM. Optionally locks or unlocks h. profile
// Return nil on success, an error otherwise.
func (bigv *bigvClient) SetVirtualMachineHardwareProfile(name VirtualMachineName, profile string, locked ...bool) (err error) {
	err = bigv.validateVirtualMachineName(&name)
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

	_, _, err = bigv.Request(true, "PUT", path, what)
	return err
}

// SetVirtualMachineHardwareProfileLock locks or unlocks the hardware profile of a VM.
// Return nil on success, an error otherwise.
func (bigv *bigvClient) SetVirtualMachineHardwareProfileLock(name VirtualMachineName, locked bool) (err error) {
	err = bigv.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	what := `{"hardware_profile_locked": false}`
	if locked {
		what = `{"hardware_profile_locked": true}`
	}

	_, _, err = bigv.Request(true, "PUT", path, what)
	return err
}
