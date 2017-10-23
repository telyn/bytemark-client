package lib_test

import (
	"encoding/json"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestCreateBackupSchedule(t *testing.T) {
	name := lib.VirtualMachineName{
		VirtualMachine: "labelmaker",
		Group:          "chocolatefactory",
		Account:        "wonka",
	}

	testSchedule := brain.BackupSchedule{
		StartDate: "00:00",
		Interval:  3306,
	}

	expectedBody, err := json.Marshal(testSchedule)
	if err != nil {
		t.Fatal(err)
	}

	rts := testutil.RequestTestSpec{
		Endpoint:      lib.BrainEndpoint,
		URL:           "/accounts/wonka/groups/chocolatefactory/virtual_machines/labelmaker/discs/disc-label/backup_schedules",
		Method:        "POST",
		AssertRequest: assert.BodyString(string(expectedBody)),
		Response:      json.RawMessage(`{"start_date": "2017-02-06T11:29:35+00:00", "interval_seconds": 3540, "id": 9}`),
	}
	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		_, err = client.CreateBackupSchedule(name, "disc-label", testSchedule.StartDate, testSchedule.Interval)
		if err != nil {
			t.Errorf("TestCreateBackupSchedule ERR: %v", err)
		}
	})
}

func TestDeleteBackupSchedule(t *testing.T) {
	name := lib.VirtualMachineName{
		VirtualMachine: "labelmaker",
		Group:          "chocolatefactory",
		Account:        "wonka",
	}

	rts := testutil.RequestTestSpec{
		Endpoint: lib.BrainEndpoint,
		URL:      "/accounts/wonka/groups/chocolatefactory/virtual_machines/labelmaker/discs/disc-label/backup_schedules/324",
		Method:   "DELETE",
	}
	rts.Run(t, testutil.Name(0), true, func(client lib.Client) {
		err := client.DeleteBackupSchedule(name, "disc-label", 324)
		if err != nil {
			t.Errorf("TestDeleteBackupSchedule ERR: %v", err)
		}
	})
}
