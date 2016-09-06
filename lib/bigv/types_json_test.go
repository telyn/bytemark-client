package bigv

import (
	"encoding/json"
	"fmt"
	"github.com/cheekybits/is"
	"net"
	"testing"
)

var (
	fixtureDiscOne = `{
    "id": 99994,
    "label": "vda",
    "size": 35840,
    "storage_grade": "sata",
    "storage_pool": "tail99-sata9",
    "type": "application/vnd.bigv.disc",
    "virtual_machine_id": 99999
}`

	fixtureDiscTwo = `{
    "id": 99995,
    "label": "vdb",
    "size": 666666,
    "storage_grade": "archive",
    "storage_pool": "tail98-sata9",
    "type": "application/vnd.bigv.disc",
    "virtual_machine_id": 99999
}`

	fixtureDiscThree = `{
    "id": 99993,
    "label": "vdc",
    "size": 333333,
    "storage_grade": "archive",
    "storage_pool": "tail97-sata9",
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

	fixtureVM = `{"autoreboot_on": true,
    "cdrom_url": null,
    "cores": 1,
    "deleted": false,
    "discs": [` + fixtureDiscOne + `,` + fixtureDiscTwo + `,` + fixtureDiscThree + `
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

func TestDiscUnmarshal(t *testing.T) {
	is := is.New(t)
	disc := new(Disc)
	err := json.Unmarshal([]byte(fixtureDiscOne), disc)

	if err != nil {
		panic(err)
	}

	is.Equal(99994, disc.ID)
	is.Equal("vda", disc.Label)
	is.Equal(35840, disc.Size)
	is.Equal("sata", disc.StorageGrade)
	is.Equal("tail99-sata9", disc.StoragePool)
	is.Equal(99999, disc.VirtualMachineID)
}

func containsIP(ips []*net.IP, ip string) bool {
	for _, i := range ips {
		if i.String() == ip {
			return true
		}
	}
	return false
}

func TestNicUnmarshal(t *testing.T) {
	is := is.New(t)

	nic := new(NetworkInterface)
	err := json.Unmarshal([]byte(fixtureNic), nic)

	if err != nil {
		panic(err)
	}

	is.Equal(99996, nic.ID)
	is.Equal(99999, nic.VirtualMachineID)
	is.Equal("ff:ff:ff:ff:ff:fe", nic.Mac)
	is.Equal(999, nic.VlanNum)

	is.Equal(true, containsIP(nic.IPs, "192.168.99.1"))
	is.Equal(true, containsIP(nic.IPs, "fe80::9999"))
	is.OK(t, nic.ExtraIPs["192.168.99.2"])
	is.Equal("192.168.99.1", nic.ExtraIPs["192.168.99.2"].String())

}

func TestVirtualMachineUnmarshal(t *testing.T) {
	is := is.New(t)

	vm := new(VirtualMachine)
	err := json.Unmarshal([]byte(fixtureVM), vm)

	if err != nil {
		fmt.Printf("%v\r\n", err)
		panic("Cannot continue")
	}

	is.Equal(true, vm.Autoreboot)
	is.Equal("", vm.CdromURL)
	is.Equal(1, vm.Cores)
	is.Equal(false, vm.Deleted)
	is.Equal(9999, vm.GroupID)
	is.Equal("virtio2013", vm.HardwareProfile)
	is.Equal(false, vm.HardwareProfileLocked)
	is.Equal("head99", vm.Head)
	is.Equal("example.notarealgroup.bytemark.uk0.bigv.io", vm.Hostname)
	is.Equal(99999, vm.ID)
	is.Equal("10.0.0.1", vm.ManagementAddress.String())
	is.Equal(2048, vm.Memory)
	is.Equal("example", vm.Name)
	is.Equal(true, vm.PowerOn)
	is.Equal("default", vm.ZoneName)

	disc := vm.Discs[0]

	is.Equal(99994, disc.ID)
	is.Equal("vda", disc.Label)
	is.Equal(35840, disc.Size)
	is.Equal("sata", disc.StorageGrade)
	is.Equal("tail99-sata9", disc.StoragePool)
	is.Equal(99999, disc.VirtualMachineID)

	disc = vm.Discs[1]

	is.Equal(99995, disc.ID)
	is.Equal("vdb", disc.Label)
	is.Equal(666666, disc.Size)
	is.Equal("archive", disc.StorageGrade)
	is.Equal("tail98-sata9", disc.StoragePool)
	is.Equal(99999, disc.VirtualMachineID)

	nic := vm.NetworkInterfaces[0]

	is.Equal(99996, nic.ID)
	is.Equal(99999, nic.VirtualMachineID)
	is.Equal("ff:ff:ff:ff:ff:fe", nic.Mac)
	is.Equal(999, nic.VlanNum)

	is.Equal(true, containsIP(nic.IPs, "192.168.99.1"))
	is.Equal(true, containsIP(nic.IPs, "fe80::9999"))
	is.OK(t, nic.ExtraIPs["192.168.99.2"])
	is.Equal("192.168.99.1", nic.ExtraIPs["192.168.99.2"].String())
}

func TestTotalDiscSize(t *testing.T) {
	is := is.New(t)

	vm := new(VirtualMachine)
	err := json.Unmarshal([]byte(fixtureVM), vm)
	if err != nil {
		panic(err)
	}

	// FixtureDiscOne + fixtureDiscTwo + fixtureDiscThree
	is.Equal(35840+666666+333333, vm.TotalDiscSize(""))
	// fixtureDiscTwo + fixtureDiscThree
	is.Equal(999999, vm.TotalDiscSize("archive"))
}

func TestImageUnmarshal(t *testing.T) {
	//is := is.New(t)
}
