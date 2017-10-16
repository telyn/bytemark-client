package lib_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestCreateBackup(t *testing.T) {
	vm := lib.VirtualMachineName{
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

	rts := testutil.RequestTestSpec{
		Method:   "POST",
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/backups",
		Response: json.RawMessage(`
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
			}`),
	}

	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		backup, err := client.CreateBackup(vm, disc)
		if err != nil {
			t.Errorf("TestCreateBackup ERR: %s", err)
		}
		if !reflect.DeepEqual(backup, testBackup) {
			t.Errorf("TestCreateBackup FAIL: expected %#v but got %#v", testBackup, backup)
		}
	})
}

func TestDeleteBackup(t *testing.T) {
	vm := lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	rts := testutil.RequestTestSpec{
		Method:        "DELETE",
		URL:           "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/backups/test-backup",
		Endpoint:      lib.BrainEndpoint,
		AssertRequest: assert.QueryValue("purge", "true"),
	}

	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		err := client.DeleteBackup(vm, disc, "test-backup")
		if err != nil {
			t.Error(err)
		}
	})
}

func TestGetBackups(t *testing.T) {
	vm := lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	rts := testutil.RequestTestSpec{
		Method:   "GET",
		URL:      "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/backups",
		Endpoint: lib.BrainEndpoint,
		Response: json.RawMessage(`[
			{
				"id": 509,
				"label": "backup-509"
			}, {
				"id": 533,
				"label": "backup-533"
			}
			]`),
	}

	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		backups, err := client.GetBackups(vm, disc)
		if err != nil {
			t.Error(err)
		}
		if len(backups) != 2 {
			t.Errorf("Wrong number of backups - %d expected, got %d", 2, len(backups))
		}
	})
}

func TestRestoreBackup(t *testing.T) {
	vm := lib.VirtualMachineName{
		VirtualMachine: "test-vm",
		Group:          "test-group",
		Account:        "test-account",
	}

	disc := "test-disc"

	result := make(map[string]bool)

	rts := testutil.RequestTestSpec{
		Method:   "PUT",
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/test-account/groups/test-group/virtual_machines/test-vm/discs/test-disc/backups/test-backup",
		AssertRequest: assert.BodyUnmarshal(&result, func(_ *testing.T, testName string) {
			if restore, ok := result["restore"]; !ok || !restore {
				t.Errorf("%s request body: restore wasn't set or was not true", testName)
			}
		}),
	}
	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		_, err := client.RestoreBackup(vm, disc, "test-backup")
		if err != nil {
			t.Error(err)
		}
	})
}
