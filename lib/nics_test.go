package lib

import "net"

func getFixtureNic() NetworkInterface {
	ip := net.IPv4(127, 0, 0, 2)
	return NetworkInterface{
		Label:            "",
		Mac:              "00:00:00:00:00",
		ID:               1,
		VlanNum:          1,
		IPs:              []*net.IP{&ip},
		ExtraIPs:         map[string]*net.IP{},
		VirtualMachineID: 1,
	}
}
