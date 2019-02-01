package lib_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func getFixtureVM() (vm brain.VirtualMachine) {
	disc := getFixtureDisc()
	nic := getFixtureNic()
	ip := net.IPv4(127, 0, 0, 1)

	return brain.VirtualMachine{
		Name:    "valid-vm",
		GroupID: 1,

		Autoreboot:            true,
		CdromURL:              "",
		Cores:                 1,
		Memory:                1,
		PowerOn:               true,
		HardwareProfile:       "fake-hardwareprofile",
		HardwareProfileLocked: false,
		ZoneName:              "default",
		Discs: []brain.Disc{
			disc,
		},
		ID:                1,
		ManagementAddress: ip,
		Deleted:           false,
		Hostname:          "valid-vm.default.account.fake-endpoint.example.com",
		Head:              "fakehead",
		NetworkInterfaces: []brain.NetworkInterface{
			nic,
		},
	}
}

func TestMoveVirtualMachine(t *testing.T) {
	testName := testutil.Name(0)

	vmMove := make(map[string]interface{})
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/old-account/groups/old-group": func(w http.ResponseWriter, r *http.Request) {
					assert.Method("GET")(t, testName, r)
					testutil.WriteJSON(t, w, brain.Account{ID: 101, Name: "old-group"})
				},
				"/accounts/old-account/groups/old-group/virtual_machines/rename-test": func(w http.ResponseWriter, r *http.Request) {
					assert.Method("PUT")(t, testName, r)
					assert.BodyUnmarshal(&vmMove, func(_t *testing.T, _testName string) {
						assert.Equal(t, testName, "new-name", vmMove["name"])
						assert.Equal(t, testName, 101.0, vmMove["group_id"])
					})(t, testName, r)
					// TODO(telyn): return...something?
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		oldName := lib.VirtualMachineName{VirtualMachine: "rename-test", Group: "old-group", Account: "old-account"}
		newName := oldName
		newName.VirtualMachine = "new-name"

		err := client.MoveVirtualMachine(oldName, newName)
		assert.Equal(t, testName, nil, err)
	})
}

func TestMoveServerGroup(t *testing.T) {
	testName := testutil.Name(0)
	vmMove := make(map[string]interface{})
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/old-account/groups/new-group": func(w http.ResponseWriter, r *http.Request) {
					assert.Method("GET")(t, testName, r)
					testutil.WriteJSON(t, w, brain.Account{ID: 105, Name: "new-group"})
				},
				"/accounts/old-account/groups/old-group/virtual_machines/group-test": func(w http.ResponseWriter, r *http.Request) {
					assert.Method("PUT")(t, testName, r)
					assert.BodyUnmarshal(&vmMove, func(_t *testing.T, _testName string) {
						assert.Equal(t, testName, "new-name", vmMove["name"])
						assert.Equal(t, testName, 105.0, vmMove["group_id"])
					})(t, testName, r)
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		oldName := lib.VirtualMachineName{VirtualMachine: "group-test", Group: "old-group", Account: "old-account"}
		newName := oldName
		newName.VirtualMachine = "new-name"
		newName.Group = "new-group"

		err := client.MoveVirtualMachine(oldName, newName)
		if err != nil {
			t.Error(err.Error())
		}
	})
}
func TestGetVirtualMachine(t *testing.T) {

	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		MuxHandlers: &testutil.MuxHandlers{
			Brain: testutil.Mux{
				"/accounts/account/groups/default/virtual_machines/valid-vm": func(w http.ResponseWriter, req *http.Request) {
					testutil.WriteJSON(t, w, getFixtureVM())
				},
				"/virtual_machines/123": func(w http.ResponseWriter, req *http.Request) {
					testutil.WriteJSON(t, w, getFixtureVM())
				},
			},
		},
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		vm, err := client.GetVirtualMachine(lib.VirtualMachineName{VirtualMachine: "", Group: "default", Account: "account"})
		assert.NotEqual(t, testName, nil, err)
		if _, ok := err.(lib.BadNameError); !ok {
			t.Fatalf("Expected BadNameError, got %T", err)
		}

		vm, err = client.GetVirtualMachine(lib.VirtualMachineName{VirtualMachine: "invalid-vm", Group: "default", Account: "account"})
		assert.NotEqual(t, testName, nil, err)

		vm, err = client.GetVirtualMachine(lib.VirtualMachineName{VirtualMachine: "valid-vm", Group: "", Account: "account"})
		assert.Equal(t, testName, nil, err)

		vm, err = client.GetVirtualMachine(lib.VirtualMachineName{VirtualMachine: "valid-vm", Group: "default", Account: "account"})
		assert.Equal(t, testName, nil, err)

		// Check that being just numeric is valid as well
		vm, err = client.GetVirtualMachine(lib.VirtualMachineName{VirtualMachine: "123"})
		assert.Equal(t, testName, nil, err)

		assert.Equal(t, testName, "127.0.0.1", vm.ManagementAddress.String())
		assert.Equal(t, testName, "127.0.0.2", vm.NetworkInterfaces[0].IPs[0].String())
	})
}

func TestCreateVirtualMachine(t *testing.T) {

	// TODO add more tests
	tests := []struct {
		Input     brain.VirtualMachineSpec
		Expect    brain.VirtualMachineSpec
		ExpectErr bool
	}{
		{
			brain.VirtualMachineSpec{
				VirtualMachine: brain.VirtualMachine{
					Name: "new-vm",
				},
				Discs: []brain.Disc{},
			},
			brain.VirtualMachineSpec{
				VirtualMachine: brain.VirtualMachine{
					Name: "new-vm",
				},
				Discs: nil,
			},
			false,
		},
		{
			brain.VirtualMachineSpec{
				VirtualMachine: brain.VirtualMachine{},
				Discs: []brain.Disc{
					{},
					{
						Size:         25600,
						StorageGrade: "archive",
						BackupSchedules: brain.BackupSchedules{{
							StartDate: "2017-09-13 16:35:00",
							Interval:  86400,
						}},
					},
					{
						Label: "gav",
					},
					{},
				},
			},
			brain.VirtualMachineSpec{
				VirtualMachine: brain.VirtualMachine{},
				Discs: []brain.Disc{
					{
						StorageGrade: "sata",
						Label:        "disc-1",
					},
					{
						Size:         25600,
						StorageGrade: "archive",
						Label:        "disc-2",
						BackupSchedules: brain.BackupSchedules{{
							StartDate: "2017-09-13 16:35:00",
							Interval:  86400,
						}},
					},
					{
						Label:        "gav",
						StorageGrade: "sata",
					},
					{
						StorageGrade: "sata",
						Label:        "disc-4",
					},
				},
			},
			false,
		},
	}

	for i, test := range tests {
		testName := testutil.Name(i)
		spec := brain.VirtualMachineSpec{}
		rts := testutil.RequestTestSpec{
			Method:   "POST",
			Endpoint: lib.BrainEndpoint,
			URL:      "/accounts/test-account/groups/test-group/vm_create",
			AssertRequest: assert.BodyUnmarshal(&spec, func(_ *testing.T, _ string) {
				assert.Equal(t, testName, test.Expect, spec)
			}),
			Response: test.Expect.VirtualMachine,
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			group := pathers.GroupName{Group: "test-group", Account: "test-account"}
			_, err := client.CreateVirtualMachine(group, test.Input)
			if err != nil && !test.ExpectErr {
				t.Fatal(err)
			}
		})
	}
}

func TestSetVirtualMachineCDROM(t *testing.T) {
	testurl := "test-cdrom-url"
	expected := map[string]interface{}{
		"cdrom_url": testurl,
	}
	unmarshalled := make(map[string]interface{})
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:   "PUT",
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/test-account/groups/test-group/virtual_machines/test-vm",
		AssertRequest: assert.BodyUnmarshal(&unmarshalled, func(_ *testing.T, _ string) {
			assert.Equal(t, testName, expected, unmarshalled)
		}),
	}
	rts.Run(t, testName, true, func(client lib.Client) {
		err := client.SetVirtualMachineCDROM(lib.VirtualMachineName{
			VirtualMachine: "test-vm",
			Group:          "test-group",
			Account:        "test-account",
		}, testurl)
		if err != nil {
			t.Fatal(err)
		}
	})
}
