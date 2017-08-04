package lib

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func TestCreateBackupSchedule(t *testing.T) {
	name := VirtualMachineName{
		VirtualMachine: "labelmaker",
		Group:          "chocolatefactory",
		Account:        "wonka",
	}

	testSchedule := brain.BackupSchedule{
		StartDate: "00:00",
		Interval:  3306,
	}

	seenRequest := false

	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/accounts/wonka/groups/chocolatefactory/virtual_machines/labelmaker/discs/disc-label/backup_schedules": func(wr http.ResponseWriter, r *http.Request) {
				seenRequest = true
				if r.Method != "POST" {
					t.Errorf("Wrong request method %s", r.Method)
				}
				var sched brain.BackupSchedule
				err := json.NewDecoder(r.Body).Decode(&sched)
				if err != nil {
					t.Errorf("Couldn't unmarshal - %v", err)
				}
				if sched.StartDate != testSchedule.StartDate {
					t.Errorf("Incorrect Start - expected %s, got %s", testSchedule.StartDate, sched.StartDate)
				}
				if sched.Interval != testSchedule.Interval {
					t.Errorf("Incorrect Interval - expected %d, got %d", testSchedule.Interval, sched.Interval)
				}
				_, err = wr.Write([]byte(`{"start_date": "2017-02-06T11:29:35+00:00", "interval_seconds": 3540, "id": 9}`))
				if err != nil {
					t.Errorf("error writing json %v", err)
				}
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
	_, err = client.CreateBackupSchedule(name, "disc-label", testSchedule.StartDate, testSchedule.Interval)
	if err != nil {
		t.Errorf("TestCreateBackupSchedule ERR: %v", err)
	}
	if !seenRequest {
		t.Errorf("TestCreateBackupSchedule never called the HTTP endpoint")
	}
}

func TestDeleteBackupSchedule(t *testing.T) {
	name := VirtualMachineName{
		VirtualMachine: "labelmaker",
		Group:          "chocolatefactory",
		Account:        "wonka",
	}
	seenRequest := false

	client, servers, err := mkTestClientAndServers(t, MuxHandlers{
		brain: Mux{
			"/accounts/wonka/groups/chocolatefactory/virtual_machines/labelmaker/discs/disc-label/backup_schedules/324": func(wr http.ResponseWriter, r *http.Request) {
				seenRequest = true
				if r.Method != "DELETE" {
					t.Errorf("Wrong request method %s", r.Method)
				}
				// TODO(telyn): implement...

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
	err = client.DeleteBackupSchedule(name, "disc-label", 324)
	if err != nil {
		t.Errorf("TestDeleteBackupSchedule ERR: %v", err)
	}
	if !seenRequest {
		t.Errorf("TestDeleteBackupSchedule never called the HTTP endpoint")
	}
}
