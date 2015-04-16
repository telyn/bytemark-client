package lib

import (
	"encoding/json"
	"fmt"
	"github.com/cheekybits/is"
	"reflect"
	"testing"
)

var (
	fixtureDiskOne = `{
    "id": 99994,
    "label": "vda",
    "size": 35840,
    "storage_grade": "sata",
    "storage_pool": "tail99-sata9",
    "type": "application/vnd.bigv.disc",
    "virtual_machine_id": 99999
}`

	fixtureDiskTwo = `{
    "id": 99995,
    "label": "vdb",
    "size": 666666,
    "storage_grade": "archive",
    "storage_pool": "tail98-sata9",
    "type": "application/vnd.bigv.disc",
    "virtual_machine_id": 99999
}`

	fixtureNic = `{
    "extra_ips": {
	"192.168.99.2": "192.168.99.1"
    },
    "id": 99996,
    "ips": [
	"192.168.99.1",
	"fe80::9999"
    ],
    "label": null,
    "mac": "ff:ff:ff:ff:ff:fe",
    "type": "application/vnd.bigv.network-interface",
    "virtual_machine_id": 99999,
    "vlan_num": 999
}`

	fixtureVm = `{"autoreboot_on": true,
    "cdrom_url": null,
    "cores": 1,
    "deleted": false,
    "discs": [` + fixtureDiskOne + `,` + fixtureDiskTwo + `
    ],
    "group_id": 9999,
    "hardware_profile": "virtio2013",
    "hardware_profile_locked": false,
    "head": "head99",
    "hostname": "example.notarealgroup.bytemark.uk0.bigv.io",
    "id": 99999,
    "keymap": null,
    "last_imaged_with": "wheezy",
    "management_address": "10.0.0.1",
    "memory": 2048,
    "name": "example",
    "network_interfaces": [` + fixtureNic + `
    ],
    "power_on": true,
    "type": "application/vnd.bigv.virtual-machine",
    "zone_name": "default"
}`
)

// Contains loops over the given list looking for element
// Returns (true, true) if found, (true, false) if not, and (false, false) if there was an error
func Contains(list, element interface{}) (found bool) {
	listValue := reflect.ValueOf(list)

	for i := 0; i < listValue.Len(); i++ {
		if listValue.Index(i).Interface() == element {
			return true
		}
	}
	return false
}

func TestDiskUnmarshal(t *testing.T) {
	is := is.New(t)
	disk := new(Disk)
	err := json.Unmarshal([]byte(fixtureDiskOne), disk)

	if err != nil {
		panic(err)
	}

	is.Equal(99994, disk.Id)
	is.Equal("vda", disk.Label)
	is.Equal(35840, disk.Size)
	is.Equal("sata", disk.StorageGrade)
	is.Equal("tail99-sata9", disk.StoragePool)
	is.Equal(99999, disk.VirtualMachineId)
}

func TestNicUnmarshal(t *testing.T) {
	is := is.New(t)

	nic := new(NetworkInterface)
	err := json.Unmarshal([]byte(fixtureNic), nic)

	if err != nil {
		panic(err)
	}

	is.Equal(99996, nic.Id)
	is.Equal(99999, nic.VirtualMachineId)
	is.Equal("ff:ff:ff:ff:ff:fe", nic.Mac)
	is.Equal(999, nic.VlanNum)

	is.Equal(true, Contains(nic.Ips, "192.168.99.1"))
	is.Equal(true, Contains(nic.Ips, "fe80::9999"))
	is.OK(t, nic.ExtraIps["192.168.99.2"])
	is.Equal("192.168.99.1", nic.ExtraIps["192.168.99.2"])

}

func TestVirtualMachineUnmarshal(t *testing.T) {
	is := is.New(t)

	vm := new(VirtualMachine)
	err := json.Unmarshal([]byte(fixtureVm), vm)

	if err != nil {
		fmt.Printf("%v\r\n", err)
		panic("Cannot continue")
	}

	is.Equal(true, vm.Autoreboot)
	is.Equal("", vm.CdromUrl)
	is.Equal(1, vm.Cores)
	is.Equal(false, vm.Deleted)
	is.Equal(9999, vm.GroupId)
	is.Equal("virtio2013", vm.HardwareProfile)
	is.Equal(false, vm.HardwareProfileLocked)
	is.Equal("head99", vm.Head)
	is.Equal("example.notarealgroup.bytemark.uk0.bigv.io", vm.Hostname)
	is.Equal(99999, vm.Id)
	is.Equal("10.0.0.1", vm.ManagementAddress)
	is.Equal(2048, vm.Memory)
	is.Equal("example", vm.Name)
	is.Equal(true, vm.PowerOn)
	is.Equal("default", vm.ZoneName)

	disk := vm.Discs[0]

	is.Equal(99994, disk.Id)
	is.Equal("vda", disk.Label)
	is.Equal(35840, disk.Size)
	is.Equal("sata", disk.StorageGrade)
	is.Equal("tail99-sata9", disk.StoragePool)
	is.Equal(99999, disk.VirtualMachineId)

	disk = vm.Discs[1]

	is.Equal(99995, disk.Id)
	is.Equal("vdb", disk.Label)
	is.Equal(666666, disk.Size)
	is.Equal("archive", disk.StorageGrade)
	is.Equal("tail98-sata9", disk.StoragePool)
	is.Equal(99999, disk.VirtualMachineId)

	nic := vm.NetworkInterfaces[0]

	is.Equal(99996, nic.Id)
	is.Equal(99999, nic.VirtualMachineId)
	is.Equal("ff:ff:ff:ff:ff:fe", nic.Mac)
	is.Equal(999, nic.VlanNum)

	is.Equal(true, Contains(nic.Ips, "192.168.99.1"))
	is.Equal(true, Contains(nic.Ips, "fe80::9999"))
	is.OK(t, nic.ExtraIps["192.168.99.2"])
	is.Equal("192.168.99.1", nic.ExtraIps["192.168.99.2"])
}

func TestImageUnmarshal(t *testing.T) {
	//is := is.New(t)
}
