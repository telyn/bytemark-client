package brain

type MigrationJobQueue struct {
	Discs []int `json:"discs,omitempty"`
}

type MigrationLocations struct {
	Discs []string `json:"discs,omitempty"`
	Pools []string `json:"pools,omitempty"`
	Tails []string `json:"tails,omitempty"`
}

type MigrationOptions struct {
	Priority int `json:"priority,omitempty"`
}

type MigrationDiscStatus struct {
	Done      []int `json:"done,omitempty"`
	Errored   []int `json:"errored,omitempty"`
	Cancelled []int `json:"cancelled,omitempty"`
	Skipped   []int `json:"skipped,omitempty"`
}

type MigrationStatus struct {
	Discs MigrationDiscStatus `json:"discs,omitempty"`
}

type MigrationJob struct {
	ID           int                `json:"id,omitempty"`
	Queue        MigrationJobQueue  `json:"queue,omitempty"`
	Status       MigrationStatus    `json:"status,omitempty"`
	Sources      MigrationLocations `json:"sources,omitempty"`
	Destinations MigrationLocations `json:"destinations,omitempty"`
	Options      MigrationOptions   `json:"options,omitempty"`
}
