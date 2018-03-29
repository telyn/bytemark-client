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

func TestGetMigrationJobs(t *testing.T) {
	testName := testutil.Name(0)

	testMigrationJobs := brain.MigrationJobs{{
		ID: 123,
		Args: brain.MigrationJobSpec{
			Sources: brain.MigrationJobLocations{
				Discs: []util.NumberOrString{"3"},
			},
			Destinations: brain.MigrationJobLocations{
				Pools: []util.NumberOrString{"pool.21"},
			},
		},
		Queue: brain.MigrationJobQueue{
			Discs: []int{3},
		},
		Destinations: brain.MigrationJobDestinations{
			Pools: []int{21},
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
	}}

	rts := testutil.RequestTestSpec{
		Method:   "GET",
		URL:      "/admin/migration_jobs",
		AssertRequest: assert.QueryValue("unfinished", "1"),
		Endpoint: lib.BrainEndpoint,
		Response: json.RawMessage(`[{
		"id": 123,
		"args": {
		    "options": {},
		    "sources": {
			"discs": [3]
		    },
		    "destinations": {
			"pools": ["pool.21"]
		    }
		},
		"queue": {
		    "discs": [3]
		},
		"destinations": {
		    "pools": [21]
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
	    }]`),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		migrationJobs, err := brainMethods.GetMigrationJobs(client)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, migrationJobs, testMigrationJobs)
	})
}
