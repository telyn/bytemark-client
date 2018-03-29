package brain

import (
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// GetMigrationJobActive Migrations returns a list of active migrations
// associated with the given migration job.
func GetMigrationJobActiveMigrations(client lib.Client, id int) (ms brain.Migrations, err error) {
	r, err := client.BuildRequest("GET", lib.BrainEndpoint, "/admin/migration_jobs/%s/migrations", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &ms)
	return
}
