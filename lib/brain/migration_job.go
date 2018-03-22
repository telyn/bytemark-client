package brain

// MigrationJobQueue is a list of disc IDs that are still to be migrated as
// part of a migration job.
type MigrationJobQueue struct {
	Discs []int `json:"discs,omitempty"`
}

// MigrationLocations represents source or target locations for a migration job.
type MigrationLocations struct {
	Discs []string `json:"discs,omitempty"`
	Pools []string `json:"pools,omitempty"`
	Tails []string `json:"tails,omitempty"`
}

// MigrationOptions represents options on a migration job.
type MigrationOptions struct {
	Priority int `json:"priority,omitempty"`
}

type MigrationDiscStatus struct {
	Done      []int `json:"done,omitempty"`
	Errored   []int `json:"errored,omitempty"`
	Cancelled []int `json:"cancelled,omitempty"`
	Skipped   []int `json:"skipped,omitempty"`
}

// MigrationStatus captures the status of a migration job, currently only
// discs.
type MigrationStatus struct {
	Discs MigrationDiscStatus `json:"discs,omitempty"`
}

// MigrationJob is a representation of a migration job.
type MigrationJob struct {
	ID           int                `json:"id,omitempty"`
	Queue        MigrationJobQueue  `json:"queue,omitempty"`
	Status       MigrationStatus    `json:"status,omitempty"`
	Sources      MigrationLocations `json:"sources,omitempty"`
	Destinations MigrationLocations `json:"destinations,omitempty"`
	Options      MigrationOptions   `json:"options,omitempty"`
}
