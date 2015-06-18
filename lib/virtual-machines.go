package lib

import "encoding/json"

var i2b = [...]string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

// GetVirtualMachine requests an overview of the named VM, regardless of its deletion status.
func (bigv *bigvClient) GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error) {
	vm = new(VirtualMachine)
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s?include_deleted=true&view=overview", name.Account, name.Group, name.VirtualMachine)

	err = bigv.RequestAndUnmarshal(true, "GET", path, "", vm)
	if err != nil {
		vm = nil
	}
	return vm, err
}

// DeleteVirtualMachine deletes the named virtual machine.
// returns nil on success or an error otherwise.
func (bigv *bigvClient) DeleteVirtualMachine(name VirtualMachineName, purge bool) (err error) {
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)
	if purge {
		path += "?purge=true"
	}

	_, _, err = bigv.Request(true, "DELETE", path, "")
	return err
}

// UndeleteVirtualMachine changes the deleted flag on a VM back to false.
// Return nil on success, an error otherwise.
func (bigv *bigvClient) UndeleteVirtualMachine(name VirtualMachineName) (err error) {
	path := BuildURL("/accounts/%s/groups/%s/virtual_machines/%s", name.Account, name.Group, name.VirtualMachine)

	_, _, err = bigv.Request(true, "PUT", path, `{"deleted":false}`)
	return err
}

//CreateVirtualMachine creates a virtual machine in the given group.
func (bigv *bigvClient) CreateVirtualMachine(group GroupName, spec VirtualMachineSpec) (vm *VirtualMachine, err error) {
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
		disc["label"] = d.Label
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
