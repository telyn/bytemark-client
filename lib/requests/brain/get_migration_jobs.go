package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// GetMigrationJobs returns an array of unfinished migration jobs
func GetMigrationJobs(client lib.Client) (mjs []brain.MigrationJob, err error) {
	r, err := client.BuildRequest("GET", lib.BrainEndpoint, "/admin/migration_jobs?unfinished=1%s", "")
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &mjs)
	return
}
