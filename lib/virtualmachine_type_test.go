package lib

import (
	"net"
	"testing"
)

func getFixtureVMWithManyIPs() (vm VirtualMachine, v4 []string, v6 []string) {
	vm = getFixtureVM()
	vm.NetworkInterfaces = make([]*NetworkInterface, 1)
	vm.NetworkInterfaces[0] = &NetworkInterface{
		Label: "test-nic",
		Mac:   "FF:FE:FF:FF:FF",
		IPs: []*net.IP{
			&net.IP{192, 168, 1, 16},
			&net.IP{192, 168, 1, 22},
			&net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x10},
			&net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x00},
		},
		ExtraIPs: map[string]*net.IP{
			"192.168.2.1":  &net.IP{192, 168, 1, 16},
			"192.168.5.34": &net.IP{192, 168, 1, 22},
			"fe80::1:1": &net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x00},
			"fe80::2:1": &net.IP{0xfe, 0x80, 0x00, 0x00,
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
		if s == false {
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
		if s == false {
			t.Error(ip + " was missing")
		}
	}
}
