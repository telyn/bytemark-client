package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
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

