package lib

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func getFixtureVMWithManyIPs() (vm brain.VirtualMachine, v4 []string, v6 []string) {
	vm = getFixtureVM()
	vm.NetworkInterfaces = make([]brain.NetworkInterface, 1)
	vm.NetworkInterfaces[0] = brain.NetworkInterface{
		Label: "test-nic",
		Mac:   "FF:FE:FF:FF:FF",
		IPs: []net.IP{
			net.IP{192, 168, 1, 16},
			net.IP{192, 168, 1, 22},
			net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x10},
			net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x00},
		},
		ExtraIPs: map[string]net.IP{
			"192.168.2.1":  net.IP{192, 168, 1, 16},
			"192.168.5.34": net.IP{192, 168, 1, 22},
			"fe80::1:1": net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x00},
			"fe80::2:1": net.IP{0xfe, 0x80, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x10},
		},
	}
	v4 = []string{"192.168.1.16", "192.168.1.22", "192.168.2.1", "192.168.5.34"}
	v6 = []string{"fe80::10", "fe80::100", "fe80::1:1", "fe80::2:1"}
	return
}
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
	is := is.New(t)

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/accounts/old-account/groups/old-group" {
				_, err := w.Write([]byte(`{"id":101, "name": "old-group"}`))
				if err != nil {
					t.Fatal(err)
				}
			} else if req.URL.Path == "/accounts/old-account/groups/old-group/virtual_machines/rename-test" {
				if req.Method == "PUT" {
					decoded := make(map[string]interface{})
					body, err := ioutil.ReadAll(req.Body)
					if err != nil {
						t.Fatal(err)
					}
					err = json.Unmarshal(body, &decoded)
					if err != nil {
						t.Fatal(err)
					}
					is.Equal("new-name", decoded["name"])
					is.Equal(101, decoded["group_id"])
					_, err = w.Write(body)
					if err != nil {
						t.Fatal(err)
					}
				}
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
		}),
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	oldName := VirtualMachineName{VirtualMachine: "rename-test", Group: "old-group", Account: "old-account"}
	newName := oldName
	newName.VirtualMachine = "new-name"

	err = client.MoveVirtualMachine(oldName, newName)
	if err != nil {
		t.Log(err.Error())
	}
	is.Nil(err)
}

func TestMoveServerGroup(t *testing.T) {
	is := is.New(t)

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/accounts/old-account/groups/new-group" {
				_, err := w.Write([]byte(`{"id":105, "name": "new-group"}`))
				if err != nil {
					t.Fatal(err)
				}
			} else if req.URL.Path == "/accounts/old-account/groups/old-group/virtual_machines/group-test" {
				if req.Method == "PUT" {
					decoded := make(map[string]interface{})
					body, err := ioutil.ReadAll(req.Body)
					if err != nil {
						t.Fatal(err)
					}
					err = json.Unmarshal(body, &decoded)
					if err != nil {
						t.Fatal(err)
					}
					is.Equal("new-name", decoded["name"])
					is.Equal(105, decoded["group_id"])
					_, err = w.Write(body)
					if err != nil {
						t.Fatal(err)
					}
				}
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
		}),
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	oldName := VirtualMachineName{VirtualMachine: "group-test", Group: "old-group", Account: "old-account"}
	newName := oldName
	newName.VirtualMachine = "new-name"
	newName.Group = "new-group"

	err = client.MoveVirtualMachine(oldName, newName)
	if err != nil {
		t.Log(err.Error())
	}
}
func TestGetVirtualMachine(t *testing.T) {
	is := is.New(t)
	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/accounts/account/groups/default/virtual_machines/valid-vm" || req.URL.Path == "/virtual_machines/123" {
				str, err := json.Marshal(getFixtureVM())
				if err != nil {
					t.Fatal(err)
				}
				_, err = w.Write(str)
				if err != nil {
					t.Fatal(err)
				}
			} else if req.URL.Path == "/accounts/account/groups/default/virtual_machines/invalid-vm" {
				http.NotFound(w, req)
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}

		}),
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	vm, err := client.GetVirtualMachine(VirtualMachineName{VirtualMachine: "", Group: "default", Account: "account"})
	is.NotNil(err)
	if _, ok := err.(BadNameError); !ok {
		t.Fatalf("Expected BadNameError, got %T", err)
	}

	vm, err = client.GetVirtualMachine(VirtualMachineName{VirtualMachine: "invalid-vm", Group: "default", Account: "account"})
	is.NotNil(err)

	vm, err = client.GetVirtualMachine(VirtualMachineName{VirtualMachine: "valid-vm", Group: "", Account: "account"})
	is.Nil(err)

	vm, err = client.GetVirtualMachine(VirtualMachineName{VirtualMachine: "valid-vm", Group: "default", Account: "account"})
	is.Nil(err)

	// Check that being just numeric is valid as well
	vm, err = client.GetVirtualMachine(VirtualMachineName{VirtualMachine: "123"})
	is.Nil(err)

	is.Equal("127.0.0.1", vm.ManagementAddress.String())
	is.Equal("127.0.0.2", vm.NetworkInterfaces[0].IPs[0].String())

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
		client, servers, err := mkTestClientAndServers(t, Handlers{
			brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path == "/accounts/test-account/groups/test-group/vm_create" {
					bytes, err := ioutil.ReadAll(req.Body)
					if err != nil {
						t.Fatalf("#%d - %v", i, err)
					}
					var spec brain.VirtualMachineSpec
					err = json.Unmarshal(bytes, &spec)

					if err != nil {
						t.Fatalf("#%d - %v", i, err)
					}
					js, err := json.MarshalIndent(spec, "IN", "    ")
					if err != nil {
						t.Fatal(err)
					}
					t.Log(string(js))
					bytes, err = json.Marshal(test.Expect.VirtualMachine)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(test.Expect, spec) {
						t.Error("spec did not deep-equal what was expected.")
					} else {
						_, err = w.Write(bytes)
						if err != nil {
							t.Fatal(err)
						}
					}
				} else {
					t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
				}

			}),
		})
		defer servers.Close()

		if err != nil {
			t.Fatal(err)
		}
		err = client.AuthWithCredentials(map[string]string{})
		if err != nil {
			t.Fatal(err)
		}

		group := GroupName{Group: "test-group", Account: "test-account"}
		_, err = client.CreateVirtualMachine(group, test.Input)
		if err != nil && !test.ExpectErr {
			t.Fatal(err)
		}
	}
}

func TestSetVirtualMachineCDROM(t *testing.T) {
	testurl := "test-cdrom-url"
	expected := map[string]interface{}{
		"cdrom_url": testurl,
	}
	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/accounts/test-account/groups/test-group/virtual_machines/test-vm" {
				bytes, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				unmarshalled := make(map[string]interface{})
				err = json.Unmarshal(bytes, &unmarshalled)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(expected, unmarshalled) {
					t.Error("unmarshalled map did not deep-equal what was expected.")
				} else {
					_, err = w.Write(bytes)
					if err != nil {
						t.Fatal(err)
					}
				}
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}

		}),
	})
	defer servers.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	err = client.SetVirtualMachineCDROM(VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}, testurl)
	if err != nil {
		t.Fatal(err)
	}
}
