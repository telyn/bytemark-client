package admin_test

import (
	"math/big"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

var defVM = lib.VirtualMachineName{Group: "default", Account: "default-account"}

func getFixtureVM() brain.VirtualMachine {
	return brain.VirtualMachine{
		Name:     "test-server",
		Hostname: "test-server.test-group",
		GroupID:  1,
	}
}

func getFixtureVLAN() brain.VLAN {
	return brain.VLAN{
		ID:        1,
		Num:       1,
		UsageType: "",
		IPRanges:  make([]brain.IPRange, 0),
	}
}

func getFixtureIPRange() brain.IPRange {
	return brain.IPRange{
		ID:        1,
		Spec:      "192.168.1.1/28",
		VLANNum:   1,
		Zones:     make([]string, 0),
		Available: big.NewInt(11),
	}
}

func getFixtureHead() brain.Head {
	return brain.Head{
		ID:    1,
		Label: "h1",
	}
}

func getFixtureTail() brain.Tail {
	return brain.Tail{
		ID:    1,
		Label: "t1",
	}
}

func getFixtureStoragePool() brain.StoragePool {
	return brain.StoragePool{
		Name:  "sata1",
		Label: "t1-sata1",
	}
}

func getFixtureDisc() brain.Disc {
	return brain.Disc{
		ID:    132,
		Label: "disc.sata-1.132",
	}
}

func getFixtureGroup() brain.Group {
	vms := make([]brain.VirtualMachine, 1)
	vms[0] = getFixtureVM()

	return brain.Group{
		Name:            "test-group",
		VirtualMachines: vms,
	}
}
