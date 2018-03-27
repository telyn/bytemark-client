package brain_test

import (
	"encoding/json"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"github.com/BytemarkHosting/bytemark-client/lib/util"
)

func TestCreateMigrationJob(t *testing.T) {
	testName := testutil.Name(0)
	testMigrationJobSpec := brain.MigrationJobSpec{
		Sources: brain.MigrationJobLocations{
			Tails: []util.NumberOrString{"1"},
		},
	}
	testMigrationJob := brain.MigrationJob{
		ID:   123,
		Args: testMigrationJobSpec,
		Queue: brain.MigrationJobQueue{
			Discs: []int{1},
		},
		Destinations: brain.MigrationJobDestinations{
			Pools: []int{1},
		},
		Status: brain.MigrationJobStatus{
			Discs: brain.MigrationJobDiscStatus{
				Done:      []int{},
				Errored:   []int{},
				Cancelled: []int{},
				Skipped:   []int{},
			},
		},
		Priority:  5,
		CreatedAt: "2018-03-15T14:23:00.579Z",
		UpdatedAt: "2018-03-15T15:28:28.244Z",
	}

	rts := testutil.RequestTestSpec{
		Method:   "POST",
		URL:      "/admin/migration_jobs",
		Endpoint: lib.BrainEndpoint,
		Response: json.RawMessage(`{
		"id": 123,
		"args": {
		    "options": {},
		    "sources": {
			"tails": [1]
		    }
		},
		"queue": {
		    "discs": [1]
		},
		"destinations": {
		    "pools": [1]
		},
		"status": {
		    "discs": {
			"done": [],
			"errored": [],
			"cancelled": [],
			"skipped": []
		    }
		},
		"priority": 5,
		"started_at": null,
		"finished_at": null,
		"created_at": "2018-03-15T14:23:00.579Z",
		"updated_at": "2018-03-15T15:28:28.244Z"
	    }`),
		AssertRequest: assert.Body(func(t *testing.T, testName string, body string) {
			var req brain.MigrationJobSpec
			err := json.Unmarshal([]byte(body), &req)
			if err != nil {
				t.Fatalf("failed to unmarshal request body: %v", err)
			}
			assert.Equal(t, testName, req, testMigrationJobSpec)
		}),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		migrationJob, err := brainMethods.CreateMigrationJob(client, testMigrationJobSpec)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, migrationJob, testMigrationJob)
	})
}
