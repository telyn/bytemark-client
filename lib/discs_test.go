package lib

func getFixtureDisc() Disc {
	return Disc{
		Label:            "",
		StorageGrade:     "sata",
		Size:             26400,
		Id:               1,
		VirtualMachineId: 1,
		StoragePool:      "fakepool",
	}
}
