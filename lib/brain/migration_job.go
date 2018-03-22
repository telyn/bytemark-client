package brain

import (
    "encoding/json"
)

// MigrationJobQueue is a list of disc IDs that are still to be migrated as
// part of a migration job.
type MigrationJobQueue struct {
	Discs []int `json:"discs,omitempty"`
}

// MigrationJobLocations represents source or target locations for a migration
// job. Discs, pools and tails maybe represented by ID number, label or UUID.
type MigrationJobLocations struct {
	Discs []json.Number `json:"discs,omitempty"`
	Pools []json.Number `json:"pools,omitempty"`
	Tails []json.Number `json:"tails,omitempty"`
}

// MigrationJobDestinations represents available desintations for a migration
// job. Unlike MigrationJobLocations, these are represented using ID number
// only.
type MigrationJobDestinations struct {
	Pools []int `json:"pools,omitempty"`
}

// MigrationJobOptions represents options on a migration job.
type MigrationJobOptions struct {
	Priority int `json:"priority,omitempty"`
}

// MigrationJobDiscStatus represents the current status of a migration job.
// Each entry is a list of disc IDs indicating the fate of discs that
// have been removed from the queue.
type MigrationJobDiscStatus struct {
	Done      []int `json:"done,omitempty"`
	Errored   []int `json:"errored,omitempty"`
	Cancelled []int `json:"cancelled,omitempty"`
	Skipped   []int `json:"skipped,omitempty"`
}

// MigrationJobStatus captures the status of a migration job, currently only
// discs.
type MigrationJobStatus struct {
	Discs MigrationJobDiscStatus `json:"discs,omitempty"`
}

// MigrationJobSpec is a specification of a migration job to be created
type MigrationJobSpec struct {
	Options      MigrationJobOptions   `json:"options,omitempty"`
	Sources      MigrationJobLocations `json:"sources,omitempty"`
	Destinations MigrationJobLocations `json:"destinations,omitempty"`
}

// MigrationJob is a representation of a migration job.
type MigrationJob struct {
	ID           int                      `json:"id,omitempty"`
	Args         MigrationJobSpec         `json:"args,omitempty"`
	Queue        MigrationJobQueue        `json:"queue,omitempty"`
	Destinations MigrationJobDestinations `json:"destinations,omitempty"`
	Status       MigrationJobStatus       `json:"status,omitempty"`
	Priority     int                      `json:"priority,omitempty"`
	StartedAt    string                   `json:"started_at,omitempty"`
	FinishedAt   string                   `json:"finished_at,omitempty"`
	CreatedAt    string                   `json:"created_at,omitempty"`
	UpdatedAt    string                   `json:"updated_at,omitempty"`
}
