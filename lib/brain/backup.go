package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"io"
)

// ColdStorageGrade is the name for the storage grade used as 'cold storage' - i.e. where backups get sent after being taken.
var ColdStorageGrade = "iceberg"

// Backup represents a single backup of a disc. Backups are taken on the same tail as the disc, and then migrated to a different storage grade immediately.
type Backup struct {
	Disc
	ParentDiscID int  `json:"parent_disc_id"`
	Manual       bool `json:"manual"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (s Backup) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "ID, Label, Size"
	}
	return "ID, Manual, Label, StorageGrade, Size, BackupCount, BackupSchedules"
}

// OnColdStorage returns true if the disc is currently on cold storage (whatever storage grade that is)
func (s Backup) OnColdStorage() bool {
	return s.StorageGrade == ColdStorageGrade
}

// PrettyPrint outputs a nicely-formatted string detailing the backup to the given writer.
func (s Backup) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	backupTpl := `
{{ define "backup_sgl" }}{{ .Label }}{{ if not .OnColdStorage }} (in progress){{ end }}{{ end }}

{{ define "backup_medium" }}{{ template "backup_sgl" . }}{{ end }}

{{ define "backup_full" }}{{ template "backup_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, backupTpl, "backup"+string(detail), s)
}

// Backups represents a collection of backups
type Backups []Backup

// PrettyPrint outputs a nicely-formatted string detailing the backup to the given writer.
func (ss Backups) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	backupsTpl := `
{{ define "backups_full" }}{{ template "backups_medium" . }}{{ end }}
{{ define "backups_medium" -}}
{{- range . -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "backups_sgl" }}{{ len . | pluralize "backup" "backups" }}{{ end }}
`
	return prettyprint.Run(wr, backupsTpl, "backups"+string(detail), ss)
}
