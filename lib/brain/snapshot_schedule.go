package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
)

type SnapshotSchedule struct {
	StartDate string
	Interval  int
}

func (sched SnapshotSchedule) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	scheduleTpl := `
{{ define "schedule_sgl" }}{{ printf "Every %d seconds starting from %s" .Interval .StartDate }}{{ end }}
{{ define "schedule_medium" }}{{ template "schedule_sgl" . }}{{ end }}
{{ define "schedule_full" }}{{ template "schedule_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, scheduleTpl, "schedule"+string(detail), sched)
}
