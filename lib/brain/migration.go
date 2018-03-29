package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Migration represents the migration of a single disc.
// The json returned from the brain can also represent virtual machine
// migrations, but we're ignoring that for now.
type Migration struct {
	ID             int    `json:"id,omitempty"`
	TailID         int    `json:"tail_id,omitempty"`
	DiscID         int    `json:"disc_id,omitempty"`
	Port           int    `json:"port,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	MigrationJobID int    `json:"migration_job_id,omitempty"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (m Migration) DefaultFields(f output.Format) string {
	return "ID, TailID, DiscID, Port, CreatedAt, UpdatedAt, MigrationJobID"
}

// PrettyPrint formats a Migration for display
func (m Migration) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{ define "migration_full" }} ▸ {{ .ID }}
     migration_job_id: {{ .MigrationJobID }}
     tail_id: {{ .TailID }}
     disc_id: {{ .DiscID }}
     port: {{ .Port }}
     created_at: {{ .CreatedAt }}
     updated_at: {{ .UpdatedAt }}
{{ end -}}{{- define "migration_sgl" -}}{{ .DiscID }}{{- end -}}`
	return prettyprint.Run(wr, template, "migration"+string(detail), m)
}
