package lib_test

import (
	"fmt"
	"math/big"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func testPostVirtualMachine(t *testing.T, endpoint string, vm *brain.VirtualMachine, testBody string, runTest func(client lib.Client) error) {
	testName := testutil.Name(0)

	getVMEndpoint := "/accounts/def-account/groups/def-group/virtual_machines/def-name"

	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				getVMEndpoint: func(wr http.ResponseWriter, r *http.Request) {
					assert.Method("GET")(t, testName, r)

					if vm != nil {
						testutil.WriteJSON(t, wr, vm)
					} else {
						fmt.Printf("%s vm was nil\n", testName)
						wr.WriteHeader(http.StatusNotFound)
					}
				},
				endpoint: func(wr http.ResponseWriter, r *http.Request) {
					assert.Method("POST")(t, testName, r)
					assert.BodyString(strings.TrimSpace(testBody))(t, testName, r)
				},
			},
		},
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		err := runTest(client)
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
	})
}

func testPutVirtualMachine(t *testing.T, endpoint string, vm *brain.VirtualMachine, testBody string, runTest func(client lib.Client) error) {
	testName := testutil.Name(0)

	getVMEndpoint := "/accounts/def-account/groups/def-group/virtual_machines/def-name"

	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				getVMEndpoint: func(wr http.ResponseWriter, r *http.Request) {
					assert.Method("GET")(t, testName, r)

					if vm != nil {
						testutil.WriteJSON(t, wr, vm)
					} else {
						wr.WriteHeader(http.StatusNotFound)
					}
				},
				endpoint: func(wr http.ResponseWriter, r *http.Request) {
					assert.Method("PUT")(t, testName, r)
					assert.BodyString(strings.TrimSpace(testBody))(t, testName, r)
				},
			},
		},
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		err := runTest(client)
		if err != nil {
			t.Fatalf("%s err %s", testName, err)
		}
	})
}

func simpleGetTest(t *testing.T, url string, testObject interface{}, runTest func(lib.Client) (interface{}, error)) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      url,
		Response: testObject,
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		object, err := runTest(client)
		if err != nil {
			t.Errorf("%s errored: %s", testName, err.Error())
		}

		assert.Equal(t, testName, testObject, object)
	})
}

func simpleDeleteTest(t *testing.T, url string, runTest func(lib.Client) error) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "DELETE",
		Endpoint: lib.BrainEndpoint,
		URL:      url,
	}
	rts.Run(t, testName, true, func(client lib.Client) {

		err := runTest(client)
		if err != nil {
			t.Errorf("%s errored: %s", testName, err.Error())
		}
	})
}

func simplePutTest(t *testing.T, url string, testBody string, runTest func(lib.Client) error) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:        "PUT",
		Endpoint:      lib.BrainEndpoint,
		URL:           url,
		AssertRequest: assert.BodyString(testBody),
	}
	rts.Run(t, testName, true, func(client lib.Client) {

		err := runTest(client)
		if err != nil {
			t.Errorf("%s errored: %s", testName, err.Error())
		}
	})
}

func simplePostTest(t *testing.T, url string, testBody string, runTest func(lib.Client) error) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:        "POST",
		Endpoint:      lib.BrainEndpoint,
		URL:           url,
		AssertRequest: assert.BodyString(testBody),
	}
	rts.Run(t, testName, true, func(client lib.Client) {

		err := runTest(client)
		if err != nil {
			t.Errorf("%s errored: %s", testName, err.Error())
		}
	})
}

func TestGetVLANS(t *testing.T) {
	testVLANs := brain.VLANs{
		{
			ID:        90210,
			Num:       123,
			UsageType: "recipes",
			IPRanges: brain.IPRanges{
				{
					ID:      1234,
					Spec:    "192.168.13.0/24",
					VLANNum: 123,
					Zones: []string{
						"test-zone",
					},
					Available: big.NewInt(200),
				},
			},
		},
	}

	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/admin/vlans",
		Response: testVLANs,
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		vlans, err := client.GetVLANs()
		if err != nil {
			t.Errorf("%s - Unexpected error %s", testName, err)
		}
		assert.Equal(t, testName, vlans, testVLANs)
	})
}

func TestGetVLAN(t *testing.T) {
	testVLAN := brain.VLAN{
		Num: 123,
	}
	simpleGetTest(t, "/admin/vlans/123", testVLAN, func(client lib.Client) (interface{}, error) {
		return client.GetVLAN(123)
	})
}

func TestGetIPRanges(t *testing.T) {
	testIPRanges := brain.IPRanges{
		{
			ID:      1234,
			Spec:    "192.168.13.0/24",
			VLANNum: 123,
			Zones: []string{
				"test-zone",
			},
			Available: big.NewInt(200),
		},
	}
	simpleGetTest(t, "/admin/ip_ranges", testIPRanges, func(client lib.Client) (interface{}, error) {
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
		Available: big.NewInt(200),
	}
	simpleGetTest(t, "/admin/ip_ranges/1234", testIPRange, func(client lib.Client) (interface{}, error) {
		return client.GetIPRange("1234")
	})
}

func TestGetIPRangeByIPRange(t *testing.T) {
	testName := testutil.Name(0)
	testIPRange := brain.IPRange{
		ID:      1234,
		Spec:    "192.168.13.0/24",
		VLANNum: 123,
		Zones: []string{
			"test-zone",
		},
		Available: big.NewInt(200),
	}

	rts := testutil.RequestTestSpec{
		Method:        "GET",
		Endpoint:      lib.BrainEndpoint,
		URL:           "/admin/ip_ranges",
		AssertRequest: assert.BodyFormValue("cidr", "192.168.13.0/24"),
		Response:      brain.IPRanges{testIPRange},
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		object, err := client.GetIPRange("192.168.13.0/24")
		if err != nil {
			t.Fatalf("TestGetIPRangeByIPRange errored: %s", err.Error())
		}

		if !reflect.DeepEqual(testIPRange, object) {
			t.Errorf("TestGetIPRangeByIPRange didn't get expected object.\r\nExpected: %#v\r\nActual:   %#v", testIPRange, object)
		}
	})
}

func TestGetHeads(t *testing.T) {
	testHeads := brain.Heads{
		{
			ID:       315,
			UUID:     "234833-2493-3423-324235",
			Label:    "test-head315",
			ZoneName: "awesomecoolguyzone",

			Architecture: "x86_64",
			// because of the way json Unmarshals net.IPs different to specifying them in this way this line is commented out
			// CCAddress:     &net.IP{214, 233, 32, 31},
			LastNote:      "melons",
			TotalMemory:   241000,
			UsageStrategy: "",
			Models:        []string{"generic", "intel"},

			FreeMemory:          123400,
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
			LastNote:      "more than a hundred years old",
			TotalMemory:   241000,
			UsageStrategy: "",
			Models:        []string{"generic", "intel"},

			FreeMemory:          234000,
			IsOnline:            true,
			UsedCores:           1,
			VirtualMachineCount: 1,
		},
	}
	simpleGetTest(t, "/admin/heads", testHeads, func(client lib.Client) (interface{}, error) {
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
		LastNote:      "more than a hundred years old",
		TotalMemory:   241000,
		UsageStrategy: "",
		Models:        []string{"generic", "intel"},

		FreeMemory:          234000,
		IsOnline:            true,
		UsedCores:           1,
		VirtualMachineCount: 1,
	}
	simpleGetTest(t, "/admin/heads/239", testHead, func(client lib.Client) (interface{}, error) {
		return client.GetHead("239")
	})
}

func TestGetTails(t *testing.T) {
	testTails := brain.Tails{
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
	simpleGetTest(t, "/admin/tails", testTails, func(client lib.Client) (interface{}, error) {
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
	simpleGetTest(t, "/admin/tails/1345", testTail, func(client lib.Client) (interface{}, error) {
		return client.GetTail("1345")
	})
}

func TestGetStoragePools(t *testing.T) {
	testStoragePools := brain.StoragePools{
		{
			Label:           "swimming-pool",
			Zone:            "frozone",
			Size:            244500,
			FreeSpace:       43355,
			AllocatedSpace:  20000,
			Discs:           4,
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
			Discs:           4,
			OvercommitRatio: 100,
			UsageStrategy:   "",
			StorageGrade:    "sata",
			Note:            "this note is a test",
		},
	}
	simpleGetTest(t, "/admin/storage_pools", testStoragePools, func(client lib.Client) (interface{}, error) {
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
		Discs:           4,
		OvercommitRatio: 100,
		UsageStrategy:   "",
		StorageGrade:    "sata",
		Note:            "this note is a test",
	}
	simpleGetTest(t, "/admin/storage_pools/useful-pool", testStoragePool, func(client lib.Client) (interface{}, error) {
		return client.GetStoragePool("useful-pool")
	})
}

func TestGetMigratingVMs(t *testing.T) {
	testVMs := brain.VirtualMachines{
		{
			Name: "coolvm",
		},
		{
			Name: "uncoolvm",
		},
	}
	simpleGetTest(t, "/admin/migrating_vms", testVMs, func(client lib.Client) (interface{}, error) {
		return client.GetMigratingVMs()
	})
}

func TestGetMigratingDiscs(t *testing.T) {
	testDiscs := brain.Discs{
		{
			ID:           123,
			Label:        "bliblbalb",
			StorageGrade: "sata",
			Size:         23456,
		},
		{
			ID:           1223,
			Label:        "blibsdfa",
			StorageGrade: "archive",
			Size:         24321,
		},
	}
	simpleGetTest(t, "/admin/migrating_discs", testDiscs, func(client lib.Client) (interface{}, error) {
		return client.GetMigratingDiscs()
	})
}

func TestGetStoppedEligibleVMs(t *testing.T) {
	testVMs := brain.VirtualMachines{
		{
			Name: "eligible-vm",
		}, {
			Name: "ultra-eligible-vm",
		},
	}
	simpleGetTest(t, "/admin/stopped_eligible_vms", testVMs, func(client lib.Client) (interface{}, error) {
		return client.GetStoppedEligibleVMs()
	})
}

func TestGetRecentVMs(t *testing.T) {
	testVMs := brain.VirtualMachines{
		{
			Name: "the-most-recent-vm",
		}, {
			Name: "slightly-less-recent-vm",
		},
	}
	simpleGetTest(t, "/admin/recent_vms", testVMs, func(client lib.Client) (interface{}, error) {
		return client.GetRecentVMs()
	})
}

func TestPostMigrateDiscWithNewStoragePool(t *testing.T) {
	simplePostTest(t, "/admin/discs/124/migrate", `{"new_pool_spec":"t6-sata1"}`, func(client lib.Client) error {
		return client.MigrateDisc(124, "t6-sata1")
	})
}

func TestPostMigrateDiscWithoutNewStoragePool(t *testing.T) {
	simplePostTest(t, "/admin/discs/123/migrate", `{}`, func(client lib.Client) error {
		return client.MigrateDisc(123, "")
	})
}

func TestPostMigrateVirtualMachineWithNewHead(t *testing.T) {
	testPostVirtualMachine(t, "/admin/vms/122/migrate", &brain.VirtualMachine{ID: 122}, `{"new_head_spec":"stg-h2"}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		return client.MigrateVirtualMachine(vmName, "stg-h2")
	})
}

func TestPostMigrateVirtualMachineWithoutHead(t *testing.T) {
	testPostVirtualMachine(t, "/admin/vms/121/migrate", &brain.VirtualMachine{ID: 121}, `{}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		return client.MigrateVirtualMachine(vmName, "")
	})
}

func TestPostMigrateVirtualMachineInvalidVirtualMachineName(t *testing.T) {
	testPostVirtualMachine(t, "/will-not-be-called", nil, `{}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		err := client.MigrateVirtualMachine(vmName, "")
		if _, ok := err.(lib.NotFoundError); ok {
			return nil
		}
		return err
	})
}

func TestPostReapVMs(t *testing.T) {
	simplePostTest(t, "/admin/reap_vms", "", func(client lib.Client) error {
		return client.ReapVMs()
	})
}

func TestDeleteVLAN(t *testing.T) {
	simpleDeleteTest(t, "/admin/vlans/123", func(client lib.Client) error {
		return client.DeleteVLAN(123)
	})
}

func TestPostAdminCreateGroup(t *testing.T) {
	simplePostTest(t, "/admin/groups", `{"account_spec":"test-account","group_name":"test-group"}`, func(client lib.Client) error {
		return client.AdminCreateGroup(lib.GroupName{Account: "test-account", Group: "test-group"}, 0)
	})
}

func TestPostAdminCreateGroupWithVLANNum(t *testing.T) {
	simplePostTest(t, "/admin/groups", `{"account_spec":"test-account","group_name":"test-group","vlan_num":12}`, func(client lib.Client) error {
		return client.AdminCreateGroup(lib.GroupName{Account: "test-account", Group: "test-group"}, 12)
	})
}

func TestPostCreateIPRange(t *testing.T) {
	simplePostTest(t, "/admin/ip_ranges", `{"ip_range":"192.168.1.1/24","vlan_num":123}`, func(client lib.Client) error {
		return client.CreateIPRange("192.168.1.1/24", 123)
	})
}

func TestPostCancelDiscMigration(t *testing.T) {
	simplePostTest(t, "/admin/discs/1234/cancel_migration", ``, func(client lib.Client) error {
		return client.CancelDiscMigration(1234)
	})
}

func TestPostCancelVMMigration(t *testing.T) {
	simplePostTest(t, "/admin/vms/1235/cancel_migration", ``, func(client lib.Client) error {
		return client.CancelVMMigration(1235)
	})
}

func TestPostEmptyStoragePool(t *testing.T) {
	simplePostTest(t, "/admin/storage_pools/pool1/empty", ``, func(client lib.Client) error {
		return client.EmptyStoragePool("pool1")
	})
}

func TestPostEmptyHead(t *testing.T) {
	simplePostTest(t, "/admin/heads/head1/empty", ``, func(client lib.Client) error {
		return client.EmptyHead("head1")
	})
}

func TestPostReifyDisc(t *testing.T) {
	simplePostTest(t, "/admin/discs/1231/reify", ``, func(client lib.Client) error {
		return client.ReifyDisc(1231)
	})
}

func TestPostApproveVM(t *testing.T) {
	testPostVirtualMachine(t, "/admin/vms/134/approve", &brain.VirtualMachine{ID: 134}, `{}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		return client.ApproveVM(vmName, false)
	})
}

func TestPostApproveVMAndPowerOn(t *testing.T) {
	testPostVirtualMachine(t, "/admin/vms/145/approve", &brain.VirtualMachine{ID: 145}, `{"power_on":true}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		return client.ApproveVM(vmName, true)
	})
}

func TestPostRejectVM(t *testing.T) {
	testPostVirtualMachine(t, "/admin/vms/139/reject", &brain.VirtualMachine{ID: 139}, `{"reason":"do not like the name"}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		return client.RejectVM(vmName, "do not like the name")
	})
}

func TestPostRegradeDisc(t *testing.T) {
	simplePostTest(t, "/admin/discs/1238/regrade", `{"new_grade":"newgrade"}`, func(client lib.Client) error {
		return client.RegradeDisc(1238, "newgrade")
	})
}

func TestPutUpdateVMMigration(t *testing.T) {
	testPutVirtualMachine(t, "/admin/vms/149/migrate", &brain.VirtualMachine{ID: 149}, `{}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}

		return client.UpdateVMMigration(vmName, nil, nil)
	})
}

func TestPutUpdateVMMigrationWithSpeed(t *testing.T) {
	testPutVirtualMachine(t, "/admin/vms/149/migrate", &brain.VirtualMachine{ID: 149}, `{"migration_speed":8500000000000}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		speed := int64(8500000000000)

		return client.UpdateVMMigration(vmName, &speed, nil)
	})
}

func TestPutUpdateVMMigrationWithDowntime(t *testing.T) {
	testPutVirtualMachine(t, "/admin/vms/149/migrate", &brain.VirtualMachine{ID: 149}, `{"migration_downtime":15}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		downtime := 15

		return client.UpdateVMMigration(vmName, nil, &downtime)
	})
}

func TestPutUpdateVMMigrationWithSpeedAndDowntime(t *testing.T) {
	testPutVirtualMachine(t, "/admin/vms/149/migrate", &brain.VirtualMachine{ID: 149}, `{"migration_downtime":15,"migration_speed":8500000000000}`, func(client lib.Client) error {
		vmName := lib.VirtualMachineName{Account: "def-account", Group: "def-group", VirtualMachine: "def-name"}
		speed := int64(8500000000000)
		downtime := 15

		return client.UpdateVMMigration(vmName, &speed, &downtime)
	})
}

func TestPostCreateUser(t *testing.T) {
	simplePostTest(t, "/admin/users", `{"priv_spec":"cluster_su","username":"user"}`, func(client lib.Client) error {
		return client.CreateUser("user", "cluster_su")
	})
}

func TestPostUpdateHead(t *testing.T) {
	simplePutTest(t, "/admin/heads/stg-h1", `{"usage_strategy":"empty"}`, func(client lib.Client) error {
		v := "empty"
		return client.UpdateHead("stg-h1", lib.UpdateHead{UsageStrategy: &v})
	})

	simplePutTest(t, "/admin/heads/stg-h1", `{"usage_strategy":null}`, func(client lib.Client) error {
		v := ""
		return client.UpdateHead("stg-h1", lib.UpdateHead{UsageStrategy: &v})
	})

	simplePutTest(t, "/admin/heads/stg-h1", `{"overcommit_ratio":150}`, func(client lib.Client) error {
		v := 150
		return client.UpdateHead("stg-h1", lib.UpdateHead{OvercommitRatio: &v})
	})

	simplePutTest(t, "/admin/heads/stg-h1", `{"label":"new-label"}`, func(client lib.Client) error {
		v := "new-label"
		return client.UpdateHead("stg-h1", lib.UpdateHead{Label: &v})
	})
}

func TestPostUpdateTail(t *testing.T) {
	simplePutTest(t, "/admin/tails/stg-t2", `{"usage_strategy":"empty"}`, func(client lib.Client) error {
		v := "empty"
		return client.UpdateTail("stg-t2", lib.UpdateTail{UsageStrategy: &v})
	})

	simplePutTest(t, "/admin/tails/stg-t2", `{"usage_strategy":null}`, func(client lib.Client) error {
		v := ""
		return client.UpdateTail("stg-t2", lib.UpdateTail{UsageStrategy: &v})
	})

	simplePutTest(t, "/admin/tails/stg-t2", `{"overcommit_ratio":125}`, func(client lib.Client) error {
		v := 125
		return client.UpdateTail("stg-t2", lib.UpdateTail{OvercommitRatio: &v})
	})

	simplePutTest(t, "/admin/tails/stg-t2", `{"label":"new-tail-label"}`, func(client lib.Client) error {
		v := "new-tail-label"
		return client.UpdateTail("stg-t2", lib.UpdateTail{Label: &v})
	})
}

func TestPostUpdateStoragePool(t *testing.T) {
	simplePutTest(t, "/admin/storage_pools/t3-sata1", `{"usage_strategy":"empty"}`, func(client lib.Client) error {
		v := "empty"
		return client.UpdateStoragePool("t3-sata1", brain.StoragePool{UsageStrategy: &v})
	})

	simplePutTest(t, "/admin/storage_pools/t3-sata1", `{"usage_strategy":null}`, func(client lib.Client) error {
		v := ""
		return client.UpdateStoragePool("t3-sata1", brain.StoragePool{UsageStrategy: &v})
	})

	simplePutTest(t, "/admin/storage_pools/t3-sata1", `{"overcommit_ratio":115}`, func(client lib.Client) error {
		v := 115
		return client.UpdateStoragePool("t3-sata1", brain.StoragePool{OvercommitRatio: &v})
	})

	simplePutTest(t, "/admin/storage_pools/t3-sata1", `{"label":"t3-sata2"}`, func(client lib.Client) error {
		v := "t3-sata2"
		return client.UpdateStoragePool("t3-sata1", brain.StoragePool{Label: &v})
	})
}
