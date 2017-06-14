package lib

import (
	"encoding/json"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestCreateBackup(t *testing.T) {

	vm := VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	testBackup := brain.Backup{
		Disc: brain.Disc{
			ID:               506,
			Label:            "philtesting-backup-20161122134250",
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
			if req.URL.Path != "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/backups" {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
			if req.Method != "POST" {
				t.Fatalf("Wrong method %s", req.Method)
			}
			_, err := w.Write([]byte(`
			{
				"id": 506,
				"label": "philtesting-backup-20161122134250",
				"manual": true,
				"parent_disc_id": 503,
				"size": 50,
				"storage_grade": "sata",
				"storage_pool": "t5-sata1",
				"type": "application/vnd.bigv.disc",
				"virtual_machine_id": 9
			}`))
			if err != nil {
				t.Fatal(err)
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
	backup, err := client.CreateBackup(vm, disc)
	if err != nil {
		t.Errorf("TestCreateBackup ERR: %s", err)
	}
	if !reflect.DeepEqual(backup, testBackup) {
		t.Errorf("TestCreateBackup FAIL: expected %#v but got %#v", testBackup, backup)
	}
}

func TestDeleteBackup(t *testing.T) {

	vm := VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/backups/test-backup" {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
			if req.Method != "DELETE" {
				t.Fatalf("Wrong method %s", req.Method)
			}
			if req.URL.Query().Get("purge") != "true" { // TODO(telyn): should really be parsing this with url.Values and checking that "purge" == "true"
				t.Errorf("Didn't incude the purge parameter")
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
	err = client.DeleteBackup(vm, disc, "test-backup")
	if err != nil {
		t.Error(err)
	}
}

func TestGetBackups(t *testing.T) {
	vm := VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/backups" {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}
			if req.Method != "GET" {
				t.Fatalf("Wrong method %s", req.Method)
			}
			_, err := w.Write([]byte(`[
			{
				"id": 509,
				"label": "backup-509"
			}, {
				"id": 533,
				"label": "backup-533"
			}
			]
			`))
			if err != nil {
				t.Fatal(err)
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
	backups, err := client.GetBackups(vm, disc)
	if err != nil {
		t.Error(err)
	}
	if len(backups) != 2 {
		t.Errorf("Wrong number of backups - %d expected, got %d", 2, len(backups))
	}
}

func TestRestoreBackup(t *testing.T) {
	vm := VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/backups/test-backup" {
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
			_, err = w.Write([]byte(``)) // TODO(telyn): in the future asking for a restore will return a backup (disc state prior to the restore), but not yet
			if err != nil {
				t.Fatal(err)
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
	_, err = client.RestoreBackup(vm, disc, "test-backup")
	if err != nil {
		t.Error(err)
	}
	// TODO(telyn): no tests for the first return value of RestoreBackup because it's always nil until we get back a backup
}
