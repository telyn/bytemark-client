package lib

import (
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestCreateSnapshot(t *testing.T) {

	vm := VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	testSnapshot := brain.Snapshot{
		Disc: brain.Disc{
			ID:               506,
			Label:            "philtesting-snapshot-20161122134250",
			Size:             50,
			StorageGrade:     "sata",
			StoragePool:      "t5-sata1",
			VirtualMachineID: 9,
		},
		ParentDiscID: 503,
		Manual:       true,
	}

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/snapshots" {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
			if req.Method != "POST" {
				t.Fatalf("Wrong method %s", req.Method)
			}
			w.Write([]byte(`
			{
				"id": 506,
				"label": "philtesting-snapshot-20161122134250",
				"manual": true,
				"parent_disc_id": 503,
				"size": 50,
				"storage_grade": "sata",
				"storage_pool": "t5-sata1",
				"type": "application/vnd.bigv.disc",
				"virtual_machine_id": 9
			}`))
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
	snapshot, err := client.CreateSnapshot(vm, disc)
	if err != nil {
		t.Errorf("TestCreateSnapshot ERR: %s", err)
	}
	if !reflect.DeepEqual(snapshot, testSnapshot) {
		t.Errorf("TestCreateSnapshot FAIL: expected %#v but got %#v", testSnapshot, snapshot)
	}
}

func TestDeleteSnapshot(t *testing.T) {

	vm := VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/snapshots/test-snapshot" {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
			if req.Method != "DELETE" {
				t.Fatalf("Wrong method %s", req.Method)
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
	err = client.DeleteSnapshot(vm, disc, "test-snapshot")
	if err != nil {
		t.Error(err)
	}
}

func TestGetSnapshots(t *testing.T) {
	vm := VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/snapshots" {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
			if req.Method != "GET" {
				t.Fatalf("Wrong method %s", req.Method)
			}
			w.Write([]byte(`[
			{
				"id": 509,
				"label": "snapshot-509"
			}, {
				"id": 533,
				"label": "snapshot-533"
			}
			]
			`))
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
	snapshots, err := client.GetSnapshots(vm, disc)
	if err != nil {
		t.Error(err)
	}
	if len(snapshots) != 2 {
		t.Errorf("Wrong number of snapshots - %d expected, got %d", 2, len(snapshots))
	}
}

func TestRestoreSnapshot(t *testing.T) {
	vm := VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/snapshots" {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
			if req.Method != "PUT" {
				t.Fatalf("Wrong method %s", req.Method)
			}
			bytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			obj := make(map[string]bool)
			err = json.Unmarshal(bytes, &obj)
			if err != nil {
				t.Fatal(err)
			}
			if restore, ok := obj["restore"]; !ok || !restore {
				t.Error("Restore not found or was not true")
			}
			w.Write([]byte(``)) // TODO(telyn): in the future asking for a restore will return a snapshot (disc state prior to the restore), but not yet
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
	_, err = client.RestoreSnapshot(vm, disc, "test-snapshot")
	if err != nil {
		t.Error(err)
	}
	// TODO(telyn): no tests for the first return value of RestoreSnapshot because it's always nil until we get back a snapshot
}
