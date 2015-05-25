package lib

func getFixtureNic() NetworkInterface {
	return NetworkInterface{
		Label:            "",
		Mac:              "00:00:00:00:00",
		Id:               1,
		VlanNum:          1,
		Ips:              []string{"127.0.0.2"},
		ExtraIps:         map[string]string{},
		VirtualMachineId: 1,
	}
}
