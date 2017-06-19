package brain

import (
	"bytes"
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"github.com/cheekybits/is"
	"net"
	"reflect"
	"testing"
)

func getFixtureNic() NetworkInterface {
	ip := net.IPv4(127, 0, 0, 2)
	return NetworkInterface{
		Label:            "",
		Mac:              "00:00:00:00:00",
		ID:               1,
		VlanNum:          1,
		IPs:              []net.IP{ip},
		ExtraIPs:         map[string]net.IP{},
		VirtualMachineID: 1,
	}
}

func getFixtureVM() (vm VirtualMachine) {
	disc := getFixtureDisc()
	nic := getFixtureNic()
	ip := net.IPv4(127, 0, 0, 1)

	return VirtualMachine{
		Name:    "valid-vm",
		GroupID: 1,

		Autoreboot:            true,
		CdromURL:              "",
		Cores:                 1,
		Memory:                1,
		PowerOn:               true,
		HardwareProfile:       "fake-hardwareprofile",
		HardwareProfileLocked: false,
		ZoneName:              "default",
		Discs: []Disc{
			disc,
		},
		ID:                1,
		ManagementAddress: ip,
		Deleted:           false,
		Hostname:          "valid-vm.default.account.fake-endpoint.example.com",
		Head:              "fakehead",
		NetworkInterfaces: []NetworkInterface{
			nic,
		},
	}
}
func getFixtureVMWithManyIPs() (vm VirtualMachine, v4 []string, v6 []string) {
	vm = getFixtureVM()
	vm.NetworkInterfaces = make([]NetworkInterface, 1)
	vm.NetworkInterfaces[0] = NetworkInterface{
		Label: "test-nic",
		Mac:   "FF:FE:FF:FF:FF",
		IPs: []net.IP{
			net.IP{192, 168, 1, 16},
			net.IP{192, 168, 1, 22},
			net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x10},
			net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x00},
		},
		ExtraIPs: map[string]net.IP{
			"192.168.2.1":  net.IP{192, 168, 1, 16},
			"192.168.5.34": net.IP{192, 168, 1, 22},
			"fe80::1:1": net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x00},
			"fe80::2:1": net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x10},
		},
	}
	v4 = []string{"192.168.1.16", "192.168.1.22", "192.168.2.1", "192.168.5.34"}
	v6 = []string{"fe80::10", "fe80::100", "fe80::1:1", "fe80::2:1"}
	return
}

func TestAllIPv4Addresses(t *testing.T) {
	vm, v4fix, _ := getFixtureVMWithManyIPs()
	v4addrs := vm.AllIPv4Addresses()
	if 4 != len(v4addrs) {
		t.Error("Too many v4 addresses?????")
	}
	seens := make(map[string]bool)
	for _, ip := range v4fix {
		seens[ip] = false
		for _, ip2 := range v4addrs {
			if ip == ip2.String() {
				seens[ip] = true
			}
		}
	}
	for ip, s := range seens {
		if !s {
			t.Error(ip + " was missing")
		}
	}

}

func TestAllIPv6Addresses(t *testing.T) {
	vm, _, v6fix := getFixtureVMWithManyIPs()
	v6addrs := vm.AllIPv6Addresses()
	if 4 != len(v6addrs) {
		t.Error("Too many v6 addresses?????")
	}
	seens := make(map[string]bool)
	for _, ip := range v6fix {
		seens[ip] = false
		for _, ip2 := range v6addrs {
			if ip == ip2.String() {
				seens[ip] = true
			}
		}
	}
	for ip, s := range seens {
		if !s {
			t.Error(ip + " was missing")
		}
	}
}

func TestNames(t *testing.T) {
	is := is.New(t)
	vm := getFixtureVM()
	is.Equal("valid-vm.default", vm.ShortName())
	is.Equal("valid-vm.default.account", vm.FullName())
}

func TestDiscLabelOffset(t *testing.T) {
	is := is.New(t)
	vm := getFixtureVM()
	vm.Discs = []Disc{
		{
			Label: "lalalalalala",
		},
		{
			Label: "disk-6",
		},
	}
	is.Equal(6, vm.GetDiscLabelOffset())

	vm.Discs = []Disc{
		{
			Label: "lalalalalala",
		},
		{
			Label: "disk-1",
		},
	}
	is.Equal(2, vm.GetDiscLabelOffset())
	vm.Discs = []Disc{
		{
			Label: "disk-2",
		},
		{
			Label: "disk-1",
		},
	}
	is.Equal(2, vm.GetDiscLabelOffset())
}

func TestVMJSON(t *testing.T) {
	// TODO add more tests as I convert SetVirtualMachine* functions
	tests := []struct {
		vm       VirtualMachine
		expected map[string]interface{}
	}{
		{
			vm: VirtualMachine{CdromURL: "testyo"},
			expected: map[string]interface{}{
				"cdrom_url": "testyo",
			},
		},
	}

	for i, test := range tests {
		js, err := json.Marshal(test.vm)
		if err != nil {
			t.Fatalf("TestVMJSON #%d json.Marshal failed: %v\r\n", i, err.Error())
		}
		unmarshalled := make(map[string]interface{})
		err = json.Unmarshal(js, &unmarshalled)
		if err != nil {
			t.Fatalf("TestVMJSON #%d json.Unmarshal failed: %v\r\n", i, err.Error())
		}
		if !reflect.DeepEqual(test.expected, unmarshalled) {
			t.Fatalf("TestVMJSON #%d failed.\r\nEXPECTED\r\n%#v\r\nACTUAL\r\n%#v", i, test.expected, unmarshalled)
		}
	}
}

func TestFormatVM(t *testing.T) {
	is := is.New(t)
	b := new(bytes.Buffer)
	vm, _, _ := getFixtureVMWithManyIPs()

	tests := []struct {
		in     VirtualMachine
		detail prettyprint.DetailLevel
		expt   string
	}{
		{
			in:     vm,
			detail: prettyprint.SingleLine,
			expt:   " ▸ valid-vm.default (powered on) in Default",
		},
		{
			in:     vm,
			detail: prettyprint.Medium,
			expt:   " ▸ valid-vm.default (powered on) in Default\n   192.168.1.16 - 1 core, 1MiB, 25GiB on 1 disc",
		},
		{
			in:     vm,
			detail: prettyprint.Full,
			expt: ` ▸ valid-vm.default (powered on) in Default
   192.168.1.16 - 1 core, 1MiB, 25GiB on 1 disc

    discs:
      •  - 25GiB, sata grade

    ips:
      • 192.168.1.16
      • 192.168.1.22
      • 192.168.2.1
      • 192.168.5.34
      • fe80::10
      • fe80::100
      • fe80::1:1
      • fe80::2:1
`,
		},
		{
			in:     VirtualMachine{},
			detail: "_discs",
			expt:   "",
		},
		{
			in:     VirtualMachine{},
			detail: "_spec",
			expt:   "   <nil> - 0 cores, 0MiB, no discs",
		},
	}

	for _, test := range tests {
		b.Truncate(0)
		err := test.in.PrettyPrint(b, test.detail)
		if err != nil {
			t.Error(err)
		}
		is.Equal(test.expt, b.String())
	}
}
