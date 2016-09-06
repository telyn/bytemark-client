package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/bigv"
	"net"
)

func getFixtureNic() bigv.NetworkInterface {
	ip := net.IPv4(127, 0, 0, 2)
	return bigv.NetworkInterface{
		Label:            "",
		Mac:              "00:00:00:00:00",
		ID:               1,
		VlanNum:          1,
		IPs:              []*net.IP{&ip},
		ExtraIPs:         map[string]*net.IP{},
		VirtualMachineID: 1,
	}
}
