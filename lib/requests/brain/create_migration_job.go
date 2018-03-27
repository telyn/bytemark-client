package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// CreateMigrationJob creates a new migration job as per the supplied
// specification, returning the newly created job on success or an error
// otherwise.
func CreateMigrationJob(client lib.Client, mjs brain.MigrationJobSpec) (mj brain.MigrationJob, err error) {
	req, err := client.BuildRequest("POST", lib.BrainEndpoint, "/admin/migration_jobs%s", "")
	if err != nil {
	    return
	}
	_, _, err = req.MarshalAndRun(mjs, &mj)
	return
}
