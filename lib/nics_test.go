package lib

func getFixtureNic() NetworkInterface {
	return NetworkInterface{
		Label:            "",
		Mac:              "00:00:00:00:00",
		ID:               1,
		VlanNum:          1,
		IPs:              []string{"127.0.0.2"},
		ExtraIPs:         map[string]string{},
		VirtualMachineID: 1,
	}
}
