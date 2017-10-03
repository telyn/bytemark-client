package brain

import (
	"bytes"
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Disc is a representation of a VM's disc.
type Disc struct {
	Label        string `json:"label,omitempty"`
	StorageGrade string `json:"storage_grade,omitempty"`
	Size         int    `json:"size,omitempty"`

	ID               int    `json:"id,omitempty"`
	VirtualMachineID int    `json:"virtual_machine_id,omitempty"`
	StoragePool      string `json:"storage_pool,omitempty"`

	BackupCount     int             `json:"backup_count,omitempty"`
	BackupSchedules BackupSchedules `json:"backup_schedules,omitempty"`
	BackupsEnabled  bool            `json:"backups_enabled,omitempty"`

	ResourceName      string `json:"resource_name,omitempty"`
	MigrationProgress int    `json:"migration_progress,omitempty"`
	MigrationEta      int    `json:"migration_eta,omitempty"`
	MigrationSpeed    int    `json:"migration_speed,omitempty"`

	NewStorageGrade string `json:"new_storage_grade,omitempty"`
	NewStoragePool  string `json:"new_storage_pool,omitempty"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (d Disc) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "ID, Label, StorageGrade, Size, BackupCount"
	}
	return "ID, Label, StorageGrade, Size, BackupCount, BackupSchedules"
}

// PrettyPrint outputs the disc nicely-formatted to the writer.
func (d Disc) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	tmpl := `{{ define "disc_sgl" }}{{ .Label }} - {{ gibtib .Size }}, {{ template "grade" . }}{{ template "hasbackups" . }}{{ end }}

{{ define "disc_medium" }}{{ template "disc_sgl" . }}
{{ if .BackupsEnabled }}{{ prettysprint .BackupSchedules "_sgl" }}{{ end }}{{ end }}

{{ define "disc_full" }}{{ template "disc_sgl" . }}
{{ if .BackupsEnabled }}{{ prettysprint .BackupSchedules "_medium" }}{{ end }}{{ end }}

{{ define "grade" }}
  {{- if ne .NewStorageGrade "" -}}
    {{- .NewStorageGrade }} grade (
    {{- if eq .StorageGrade "iceberg" -}} restore {{- else -}} migration {{- end }} in progress)
  {{- else -}}
    {{- .StorageGrade }} grade
  {{- end -}}
{{- end }}

{{ define "hasbackups" }}{{ if gt .BackupCount 0 }} (has {{ pluralize "backup" "backups" .BackupCount }}){{ end }}{{ end }}
`
	return prettyprint.Run(wr, tmpl, "disc"+string(detail), d)
}

// EstimateBackupScheduleSize returns an estimate for the maximum amount of iceberg
// storage this disk will use for its backups, in MiB
func (d Disc) EstimateBackupScheduleSize() int {
	totalBackups := 0
	for _, bs := range d.BackupSchedules {
		totalBackups += bs.Capacity
	}
	return d.Size * totalBackups
}

// String returns the disc formatted as a string (the same as PrettyPrint with prettyprint.SingleLine detail)
func (d Disc) String() string {
	buf := new(bytes.Buffer)
	_ = d.PrettyPrint(buf, prettyprint.SingleLine)
	return buf.String()
}

// Validate makes sure the disc has a storage grade.
func (d Disc) Validate() (*Disc, error) {
	if d.StorageGrade == "" {
		d.StorageGrade = "sata"
	}
	return &d, nil
}

// Discs represents multiple Disc objects in an output.Outputtable form.
type Discs []Disc

// DefaultFields provides the default fields for Discs, which is the same as those of Disc.
func (ds Discs) DefaultFields(f output.Format) string {
	return (Disc{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the discs to wr at the given detail level.
func (ds Discs) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	discsTpl := `
{{ define "discs_sgl" }}{{ len . }} discs{{ end }}

{{ define "discs_medium" -}}
{{- range . -}}
{{ prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "discs_full" -}}{{ template "discs_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, discsTpl, "discs"+string(detail), ds)
}
