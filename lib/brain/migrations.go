package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Migrations represents more than one migration in output.Outputtable form.
type Migrations []Migration

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type, which is the same as Migration.DefaultFields.
func (ms Migrations) DefaultFields(f output.Format) string {
	return (Migration{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the migrations to writer at the given detail level.
func (ms Migrations) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	migrationsTpl := `
{{ define "migrations_sgl" }}{{ len . }} servers{{ end }}

{{ define "migrations_medium" -}}
{{- range . -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "migrations_full" }}{{ template "migrations_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, migrationsTpl, "migrations"+string(detail), ms)
}
