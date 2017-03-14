package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type simpleGetTestFn func(Client) (interface{}, error)

func simpleGetTest(t *testing.T, url string, testObject interface{}, runTest simpleGetTestFn) {
	callerPC, _, _, _ := runtime.Caller(1)
	testName := runtime.FuncForPC(callerPC).Name()

	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			url: func(wr http.ResponseWriter, r *http.Request) {
				assertMethod(t, r, "GET")
				writeJSON(t, wr, testObject)
			},
		},
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	object, err := runTest(client)
	if err != nil {
		t.Errorf("%s errored: %s", testName, err.Error())
	}

	if !reflect.DeepEqual(testObject, object) {
		t.Errorf("%s didn't get expected object.\r\nExpected: %#v\r\nActual:   %#v", testName, testObject, object)
	}
}

type simplePostTestFn func(Client) error

func simplePostTest(t *testing.T, url string, testBody string, runTest simplePostTestFn) {
	is := is.New(t)

	callerPC, _, _, _ := runtime.Caller(1)
	testName := runtime.FuncForPC(callerPC).Name()

	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			url: func(wr http.ResponseWriter, r *http.Request) {
				assertMethod(t, r, "POST")

				body, err := ioutil.ReadAll(r.Body)
				is.Nil(err)
				is.Nil(r.Body.Close())

				is.Equal(strings.TrimSpace(string(body)), strings.TrimSpace(testBody))
			},
		},
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	err = runTest(client)
	if err != nil {
		t.Errorf("%s errored: %s", testName, err.Error())
	}
}

func TestGetVLANS(t *testing.T) {
	testVLANs := []*brain.VLAN{
		{
			ID:        90210,
			Num:       123,
			UsageType: "recipes",
			IPRanges: []*brain.IPRange{
				{
					ID:      1234,
					Spec:    "192.168.13.0/24",
					VLANNum: 123,
					Zones: []string{
						"test-zone",
					},
					Available: 200.0,
				},
			},
		},
	}
	simpleGetTest(t, "/admin/vlans", testVLANs, func(client Client) (interface{}, error) {
		return client.GetVLANs()
	})
}

func TestGetVLAN(t *testing.T) {
	testVLAN := brain.VLAN{
		Num: 123,
	}
	simpleGetTest(t, "/admin/vlans/123", &testVLAN, func(client Client) (interface{}, error) {
		return client.GetVLAN(123)
	})
}

func TestGetIPRanges(t *testing.T) {
	testIPRanges := []*brain.IPRange{
		{
			ID:      1234,
			Spec:    "192.168.13.0/24",
			VLANNum: 123,
			Zones: []string{
				"test-zone",
			},
			Available: 200.0,
		},
	}
	simpleGetTest(t, "/admin/ip_ranges", testIPRanges, func(client Client) (interface{}, error) {
		return client.GetIPRanges()
	})

}

func TestGetIPRange(t *testing.T) {
	testIPRange := brain.IPRange{
		ID:      1234,
		Spec:    "192.168.13.0/24",
		VLANNum: 123,
		Zones: []string{
			"test-zone",
		},
		Available: 200.0,
	}
	simpleGetTest(t, "/admin/ip_ranges/1234", &testIPRange, func(client Client) (interface{}, error) {
		return client.GetIPRange(1234)
	})
}

func TestGetHeads(t *testing.T) {
	testHeads := []*brain.Head{
		{
			ID:       315,
			UUID:     "234833-2493-3423-324235",
			Label:    "test-head315",
			ZoneName: "awesomecoolguyzone",

			Architecture: "x86_64",
			// because of the way json Unmarshals net.IPs different to specifying them in this way this line is commented out
			// CCAddress:     &net.IP{214, 233, 32, 31},
			Note:          "melons",
			Memory:        241000,
			UsageStrategy: "",
			Models:        []string{"generic", "intel"},

			MemoryFree:          123400,
			IsOnline:            true,
			UsedCores:           9,
			VirtualMachineCount: 3,
		}, {
			ID:       239,
			UUID:     "235670-2493-3423-324235",
			Label:    "test-head239",
			ZoneName: "awesomecoolguyzone",

			Architecture: "x86_64",
			// because of the way json Unmarshals net.IPs different to specifying them in this way this line is commented out
			// CCAddress:     &net.IP{24, 43, 32, 49},
			Note:          "more than a hundred years old",
			Memory:        241000,
			UsageStrategy: "",
			Models:        []string{"generic", "intel"},

			MemoryFree:          234000,
			IsOnline:            true,
			UsedCores:           1,
			VirtualMachineCount: 1,
		},
	}
	simpleGetTest(t, "/admin/heads", testHeads, func(client Client) (interface{}, error) {
		return client.GetHeads()
	})

}

func TestGetHead(t *testing.T) {
	testHead := brain.Head{
		ID:       239,
		UUID:     "235670-2493-3423-324235",
		Label:    "test-head239",
		ZoneName: "awesomecoolguyzone",

		Architecture: "x86_64",
		// because of the way json Unmarshals net.IPs different to specifying them in this way this line is commented out
		// CCAddress:     &net.IP{24, 43, 32, 49},
		Note:          "more than a hundred years old",
		Memory:        241000,
		UsageStrategy: "",
		Models:        []string{"generic", "intel"},

		MemoryFree:          234000,
		IsOnline:            true,
		UsedCores:           1,
		VirtualMachineCount: 1,
	}
	simpleGetTest(t, "/admin/heads/239", &testHead, func(client Client) (interface{}, error) {
		return client.GetHead("239")
	})
}

func TestGetTails(t *testing.T) {
	testTails := []*brain.Tail{
		{
			ID:           1345,
			UUID:         "idont-reallyknowwhat-uuids-looklike",
			Label:        "coolTailForCoolDiscs",
			ZoneName:     "frozone",
			IsOnline:     false,
			StoragePools: []string{"swimming", "paddling"},
		}, {
			ID:           1235,
			UUID:         "888888-8888-8888-888888",
			Label:        "eight",
			ZoneName:     "eighth zone",
			IsOnline:     true,
			StoragePools: []string{"pool-eight"},
		},
	}
	simpleGetTest(t, "/admin/tails", testTails, func(client Client) (interface{}, error) {
		return client.GetTails()
	})
}

func TestGetTail(t *testing.T) {
	testTail := brain.Tail{
		ID:           1345,
		UUID:         "idont-reallyknowwhat-uuids-looklike",
		Label:        "coolTailForCoolDiscs",
		ZoneName:     "frozone",
		IsOnline:     false,
		StoragePools: []string{"swimming", "paddling"},
	}
	simpleGetTest(t, "/admin/tails/1345", &testTail, func(client Client) (interface{}, error) {
		return client.GetTail("1345")
	})
}

func TestGetStoragePools(t *testing.T) {
	testStoragePools := []*brain.StoragePool{
		{
			Label:           "swimming-pool",
			Zone:            "frozone",
			Size:            244500,
			FreeSpace:       43355,
			AllocatedSpace:  20000,
			Discs:           []string{"daves-disc", "steves-disc", "philomenas-disc", "celestes-disc"},
			OvercommitRatio: 9000,
			UsageStrategy:   "arsony",
			StorageGrade:    "wet",
			Note:            "probably best to avoid using this one",
		}, {
			Label:           "useful-pool",
			Zone:            "serious-zone",
			Size:            244500000,
			FreeSpace:       43355000,
			AllocatedSpace:  20000000,
			Discs:           []string{"daves-disc", "steves-disc", "philomenas-disc", "celestes-disc"},
			OvercommitRatio: 100,
			UsageStrategy:   "",
			StorageGrade:    "sata",
			Note:            "this note is a test",
		},
	}
	simpleGetTest(t, "/admin/storage_pools", testStoragePools, func(client Client) (interface{}, error) {
		return client.GetStoragePools()
	})
}

func TestGetStoragePool(t *testing.T) {
	testStoragePool := brain.StoragePool{
		Label:           "useful-pool",
		Zone:            "serious-zone",
		Size:            244500000,
		FreeSpace:       43355000,
		AllocatedSpace:  20000000,
		Discs:           []string{"daves-disc", "steves-disc", "philomenas-disc", "celestes-disc"},
		OvercommitRatio: 100,
		UsageStrategy:   "",
		StorageGrade:    "sata",
		Note:            "this note is a test",
	}
	simpleGetTest(t, "/admin/storage_pools/useful-pool", &testStoragePool, func(client Client) (interface{}, error) {
		return client.GetStoragePool("useful-pool")
	})
}

func TestGetMigratingVMs(t *testing.T) {
	testVMs := []*brain.VirtualMachine{
		&brain.VirtualMachine{
			Name: "coolvm",
		},
		&brain.VirtualMachine{
			Name: "uncoolvm",
		},
	}
	simpleGetTest(t, "/admin/migrating_vms", testVMs, func(client Client) (interface{}, error) {
		return client.GetMigratingVMs()
	})
}

func TestGetStoppedEligibleVMs(t *testing.T) {
	testVMs := []*brain.VirtualMachine{
		{
			Name: "eligible-vm",
		}, {
			Name: "ultra-eligible-vm",
		},
	}
	simpleGetTest(t, "/admin/stopped_eligible_vms", testVMs, func(client Client) (interface{}, error) {
		return client.GetStoppedEligibleVMs()
	})
}

func TestGetRecentVMs(t *testing.T) {
	testVMs := []*brain.VirtualMachine{
		{
			Name: "the-most-recent-vm",
		}, {
			Name: "slightly-less-recent-vm",
		},
	}
	simpleGetTest(t, "/admin/recent_vms", testVMs, func(client Client) (interface{}, error) {
		return client.GetRecentVMs()
	})
}

func TestPostMigrateDiscWithNewStoragePool(t *testing.T) {
	simplePostTest(t, "/admin/discs/124/migrate", `{"new_pool_spec":"t6-sata1"}`, func(client Client) error {
		return client.MigrateDisc(124, "t6-sata1")
	})
}

func TestPostMigrateDiscWithoutNewStoragePool(t *testing.T) {
	simplePostTest(t, "/admin/discs/123/migrate", `{}`, func(client Client) error {
		return client.MigrateDisc(123, "")
	})
}

func TestPostMigrateVirtualMachineWithNewHead(t *testing.T) {
	simplePostTest(t, "/admin/vms/122/migrate", `{"new_head_spec":"stg-h2"}`, func(client Client) error {
		return client.MigrateVirtualMachine(122, "stg-h2")
	})
}

func TestPostMigrateVirtualMachineWithoutHead(t *testing.T) {
	simplePostTest(t, "/admin/vms/121/migrate", `{}`, func(client Client) error {
		return client.MigrateVirtualMachine(121, "")
	})
}
