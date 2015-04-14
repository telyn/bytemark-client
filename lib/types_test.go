package lib

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestDiskUnmarshal(t *testing.T) {
	disk := new(Disk)
	err := json.Unmarshal([]byte(fixtureDiskOne), disk)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, 99994, disk.Id)
	assert.Equal(t, "vda", disk.Label)
	assert.Equal(t, 35840, disk.Size)
	assert.Equal(t, "sata", disk.StorageGrade)
	assert.Equal(t, "tail99-sata9", disk.StoragePool)
	assert.Equal(t, 99999, disk.VirtualMachineId)
}

func TestNicUnmarshal(t *testing.T) {
	nic := new(NetworkInterface)
	err := json.Unmarshal([]byte(fixtureNic), nic)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, 99996, nic.Id)
	assert.Equal(t, 99999, nic.VirtualMachineId)
	assert.Equal(t, "ff:ff:ff:ff:ff:fe", nic.Mac)
	assert.Equal(t, 999, nic.VlanNum)

	assert.Contains(t, nic.Ips, "192.168.99.1")
	assert.Contains(t, nic.Ips, "fe80::9999")
	assert.NotNil(t, nic.ExtraIps["192.168.99.2"])
	assert.Equal(t, "192.168.99.1", nic.ExtraIps["192.168.99.2"])

}

func TestVirtualMachineUnmarshal(t *testing.T) {

	vm := new(VirtualMachine)
	err := json.Unmarshal([]byte(fixtureVm), vm)

	if err != nil {
		fmt.Printf("%v\r\n", err)
		panic("Cannot continue")
	}

	assert.True(t, vm.Autoreboot)
	assert.Equal(t, "", vm.CdromUrl)
	assert.Equal(t, 1, vm.Cores)
	assert.False(t, vm.Deleted)
	assert.Equal(t, 9999, vm.GroupId)
	assert.Equal(t, "virtio2013", vm.HardwareProfile)
	assert.False(t, vm.HardwareProfileLocked)
	assert.Equal(t, "head99", vm.Head)
	assert.Equal(t, "example.notarealgroup.bytemark.uk0.bigv.io", vm.Hostname)
	assert.Equal(t, 99999, vm.Id)
	assert.Equal(t, "10.0.0.1", vm.ManagementAddress)
	assert.Equal(t, 2048, vm.Memory)
	assert.Equal(t, "example", vm.Name)
	assert.True(t, vm.PowerOn)
	assert.Equal(t, "default", vm.ZoneName)

	disk := vm.Discs[0]

	assert.Equal(t, 99994, disk.Id)
	assert.Equal(t, "vda", disk.Label)
	assert.Equal(t, 35840, disk.Size)
	assert.Equal(t, "sata", disk.StorageGrade)
	assert.Equal(t, "tail99-sata9", disk.StoragePool)
	assert.Equal(t, 99999, disk.VirtualMachineId)

	disk = vm.Discs[1]

	assert.Equal(t, 99995, disk.Id)
	assert.Equal(t, "vdb", disk.Label)
	assert.Equal(t, 666666, disk.Size)
	assert.Equal(t, "archive", disk.StorageGrade)
	assert.Equal(t, "tail98-sata9", disk.StoragePool)
	assert.Equal(t, 99999, disk.VirtualMachineId)

	nic := vm.NetworkInterfaces[0]

	assert.Equal(t, 99996, nic.Id)
	assert.Equal(t, 99999, nic.VirtualMachineId)
	assert.Equal(t, "ff:ff:ff:ff:ff:fe", nic.Mac)
	assert.Equal(t, 999, nic.VlanNum)

	assert.Contains(t, nic.Ips, "192.168.99.1")
	assert.Contains(t, nic.Ips, "fe80::9999")
	assert.NotNil(t, nic.ExtraIps["192.168.99.2"])
	assert.Equal(t, "192.168.99.1", nic.ExtraIps["192.168.99.2"])
}

func TestImageUnmarshal(t *testing.T) {
}
