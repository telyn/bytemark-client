package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"net"
)

func getFixtureNic() brain.NetworkInterface {
	ip := net.IPv4(127, 0, 0, 2)
	return brain.NetworkInterface{
		Label:            "",
		Mac:              "00:00:00:00:00",
		ID:               1,
		VlanNum:          1,
		IPs:              []net.IP{ip},
		ExtraIPs:         map[string]net.IP{},
		VirtualMachineID: 1,
	}
}
