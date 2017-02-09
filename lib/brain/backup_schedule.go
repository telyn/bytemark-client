package brain

import (
	"bytes"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
	"text/template"
)

// BackupSchedule represents a schedule to take backups on. It is represented as a start date in YYYY-MM-DD hh:mm:ss format (and assuming UK timezones of some kind.)
type BackupSchedule struct {
	ID        int    `json:"id,omitempty"`
	StartDate string `json:"start_at"`
	Interval  int    `json:"interval_seconds"`
}

// PrettyPrint outputs a nicely-formatted human-readable version of the schedule to the given writer.
// All the detail levels are the same.
func (sched BackupSchedule) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	// TODO(telyn): really ought to be nicer.
	scheduleTpl := `{{ define "schedule_sgl" }}{{ printf "Every %d seconds starting from %s" .Interval .StartDate }}{{ end }}
{{ define "schedule_medium" }}{{ template "schedule_sgl" . }}{{ end }}
{{ define "schedule_full" }}{{ template "schedule_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, scheduleTpl, "schedule"+string(detail), sched)
}

// BackupSchedules represents multiple backup schedules
type BackupSchedules []*BackupSchedule

// MapTemplateFragment takes a template fragment (as if it was starting within a {{ }}) and executes it against every schedule in scheds, returning all the results as a slice of strings, or an error if one occurred.
// Is this the most heinous thing ever?
func (scheds BackupSchedules) MapTemplateFragment(templateFrag string) (strs []string, err error) {
	strs = make([]string, len(scheds))
	tmpl, err := template.New("backupschedules_maptemplatefragment").Parse(`{{` + templateFrag + `}}`)
	if err != nil {
		return
	}
	for i, s := range scheds {
		var buf bytes.Buffer

		err = tmpl.Execute(&buf, s)
		if err != nil {
			return
		}
		strs[i] = buf.String()
	}
	return
}

// PrettyPrint outputs a nicely-formatted human-readable version of the schedules to the given writer.
// detail levels:
// SingleLine - outputs one line "Backups are taken every m, n, o, & p seconds" or "No backups scheduled"
// Medium - same
// Full - outputs one line per schedule, "• #ID - " followed by the SingleLine PrettyPrint of the schedule
func (scheds BackupSchedules) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	tmpl := `{{ define "backupschedules_sgl" }}{{ if len . | ne 0 }}Backups are taken every {{ map . ".Interval" | joinWithSpecialLast ", " " & "}} seconds{{ else }}No backups scheduled{{ end }}{{ end }}
{{ define "backupschedules_medium" }}{{ template "backupschedules_sgl" .}}{{ end }}
{{ define "backupschedules_full" }}{{ range . -}}
• #{{.ID}} - {{ prettysprint . "_sgl" }}
{{ end }}{{ end }}`
	return prettyprint.Run(wr, tmpl, "backupschedules"+string(detail), scheds)
}
