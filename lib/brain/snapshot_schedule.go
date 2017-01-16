package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
)

// SnapshotSchedule represents a schedule to take snapshots on. It is represented as a start date in YYYY-MM-DD hh:mm:ss format (and assuming UK timezones of some kind.)
type SnapshotSchedule struct {
	StartDate string
	Interval  int
}

// PrettyPrint outputs a nicely-formatted human-readable version of the schedule to the given writer.
// All the detail levels are the same.
func (sched SnapshotSchedule) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	scheduleTpl := `
{{ define "schedule_sgl" }}{{ printf "Every %d seconds starting from %s" .Interval .StartDate }}{{ end }}
{{ define "schedule_medium" }}{{ template "schedule_sgl" . }}{{ end }}
{{ define "schedule_full" }}{{ template "schedule_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, scheduleTpl, "schedule"+string(detail), sched)
}
