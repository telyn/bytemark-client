package lib

import (
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// GetMigrationJob returns a single migration job, given its ID
func (c *bytemarkClient) GetMigrationJob(id int) (mj brain.MigrationJob, err error) {
	r, err := c.BuildRequest("GET", BrainEndpoint, "/admin/migration_jobs/%s", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, &mj)
	return
}
