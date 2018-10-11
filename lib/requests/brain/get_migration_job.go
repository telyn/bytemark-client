package brain

import (
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// GetMigrationJob returns a single migration job, given its ID
func GetMigrationJob(client lib.Client, id int) (mj brain.MigrationJob, err error) {
	r, err := client.BuildRequest("GET", lib.BrainEndpoint, "/admin/migration_jobs/%s", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &mj)
	return
}
