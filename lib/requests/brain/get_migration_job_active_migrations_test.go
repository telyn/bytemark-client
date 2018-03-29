package brain_test

import (
	"encoding/json"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestGetMigrationJobActiveMigrations(t *testing.T) {
	testName := testutil.Name(0)

	testMigrations := brain.Migrations{{
		ID:             1,
		DiscID:         10,
		TailID:         2,
		MigrationJobID: 123,
		CreatedAt:      "2018-03-29T08:54:50.198Z",
		UpdatedAt:      "2018-03-29T08:54:50.198Z",
	}}

	rts := testutil.RequestTestSpec{
		Method:   "GET",
		URL:      "/admin/migration_jobs/123/migrations",
		Endpoint: lib.BrainEndpoint,
		Response: json.RawMessage(`[{
		"id": 1,
		"disc_id": 10,
		"tail_id": 2,
		"migration_job_id": 123,
		"created_at": "2018-03-29T08:54:50.198Z",
		"updated_at": "2018-03-29T08:54:50.198Z"
	    }]`),
	}

	rts.Run(t, testName, true, func(client lib.Client) {
		migrations, err := brainMethods.GetMigrationJobActiveMigrations(client, 123)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testName, migrations, testMigrations)
	})
}
