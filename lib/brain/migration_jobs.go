package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// MigrationJobs represents more than one account in output.Outputtable form.
type MigrationJobs []MigrationJob

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type, which is the same as MigrationJob.DefaultFields.
func (mjs MigrationJobs) DefaultFields(f output.Format) string {
	return (MigrationJob{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the virtual machines to writer at the given detail level.
func (mjs MigrationJobs) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	migrationjobsTpl := `
{{ define "migrationjobs_sgl" }}{{ len . }} servers{{ end }}

{{ define "migrationjobs_medium" -}}
{{- range . -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "migrationjobs_full" }}{{ template "migrationjobs_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, migrationjobsTpl, "migrationjobs"+string(detail), mjs)
}
