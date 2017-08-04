package lib

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
)

//CreateVirtualMachine creates a virtual machine in the given group.
func (c *bytemarkClient) CreateVirtualMachine(group GroupName, spec brain.VirtualMachineSpec) (vm brain.VirtualMachine, err error) {
	err = c.validateGroupName(&group)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/vm_create", group.Account, group.Group)
	if err != nil {
		return
	}
	if spec.IPs != nil {
		if spec.IPs.IPv4 == "" && spec.IPs.IPv6 == "" {
			spec.IPs = nil
		}
	}
	if spec.Discs != nil {
		if len(spec.Discs) == 0 {
			spec.Discs = nil
		}
		for i, disc := range spec.Discs {
			newDisc, discErr := disc.Validate()
			if discErr != nil {
				return vm, discErr
			}
			spec.Discs[i] = *newDisc
		}
		labelDiscs(spec.Discs, 0)
	}

	oldfile := log.LogFile
	log.LogFile = nil
	_, _, err = r.MarshalAndRun(spec, &vm)
	log.LogFile = oldfile
	return vm, err
}

// DeleteVirtualMachine deletes the named virtual machine.
// returns nil on success or an error otherwise.
func (c *bytemarkClient) DeleteVirtualMachine(name VirtualMachineName, purge bool) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) GetVirtualMachine(name VirtualMachineName) (vm brain.VirtualMachine, err error) {
	var r *Request

	// If the VM name is numeric, it means it is an internal Bytemark ID,
	// so we should use a different endpoint
	if _, nErr := strconv.Atoi(name.VirtualMachine); nErr == nil {
		r, err = c.BuildRequest("GET", BrainEndpoint, "/virtual_machines/%s?include_deleted=true&view=overview", name.VirtualMachine)
	} else {
		err = c.validateVirtualMachineName(&name)
		if err != nil {
			return
		}
		r, err = c.BuildRequest("GET", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s?include_deleted=true&view=overview", name.Account, name.Group, name.VirtualMachine)
	}
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &vm)
	if err != nil {
		return
	}
	return
}

//MoveVirtualMachine moves the virtual machine to the given name, across groups if needed.
func (c *bytemarkClient) MoveVirtualMachine(oldName VirtualMachineName, newName VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&oldName)
	if err != nil {
		return
	}
	err = c.validateVirtualMachineName(&newName)
	if err != nil {
		return
	}

	// create the change we want to see in the server
	change := brain.VirtualMachine{Name: newName.VirtualMachine}
	if newName.Group != "" || newName.Account != "" {
		// get group
		groupName := GroupName{Group: newName.Group, Account: newName.Account}
		group, groupErr := c.GetGroup(groupName)
		if groupErr != nil {
			return groupErr
		}
		change.GroupID = group.ID
	}

	// PUT the change
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", oldName.Account, oldName.Group, oldName.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.MarshalAndRun(change, nil)
	return err

}

// ReimageVirtualMachine reimages the named virtual machine. This will wipe everything on the first disk in the vm and install a new OS on top of it.
// Note that the machine in question must already be powered off. Once complete, according to the API docs, the vm will be powered on but its autoreboot_on will be false.
func (c *bytemarkClient) ReimageVirtualMachine(name VirtualMachineName, image brain.ImageInstall) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/reimage", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	oldfile := log.LogFile
	log.LogFile = nil
	_, _, err = r.MarshalAndRun(image, nil)
	log.LogFile = oldfile
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
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/signal", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.Run(bytes.NewBufferString(`{"signal":"reset"}`), nil)
	return err
}

// RestartVirtualMachine restarts the named virtual machine. This is
// returns nil on success or an error otherwise.
func (c *bytemarkClient) RestartVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) StartVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) StopVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) ShutdownVirtualMachine(name VirtualMachineName, stayoff bool) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) UndeleteVirtualMachine(name VirtualMachineName) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) SetVirtualMachineHardwareProfile(name VirtualMachineName, profile string, locked ...bool) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) SetVirtualMachineHardwareProfileLock(name VirtualMachineName, locked bool) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) SetVirtualMachineMemory(name VirtualMachineName, memory int) (err error) {
	err = c.validateVirtualMachineName(&name)
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
func (c *bytemarkClient) SetVirtualMachineCores(name VirtualMachineName, cores int) (err error) {
	err = c.validateVirtualMachineName(&name)
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

// SetVirtualMachineCDROM sets the URL of a CD to attach to a virtual machine. Set url to "" to remove the CD.
// Returns nil on success, an error otherwise.
func (c *bytemarkClient) SetVirtualMachineCDROM(name VirtualMachineName, url string) (err error) {
	err = c.validateVirtualMachineName(&name)
	if err != nil {
		return err
	}
	r, err := c.BuildRequest("PUT", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if err != nil {
		return err
	}

	_, _, err = r.MarshalAndRun(brain.VirtualMachine{CdromURL: url}, nil)
	return err
}
