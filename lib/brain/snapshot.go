package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
)

// ColdStorageGrade is the name for the storage grade used as 'cold storage' - i.e. where snapshots get sent after being taken.
var ColdStorageGrade = "iceberg"

// Snapshot represents a single snapshot of a disc. Snapshots are taken on the same tail as the disc, and then migrated to a different storage grade immediately.
type Snapshot struct {
	Disc
	ParentDiscID int  `json:"parent_disc_id"`
	Manual       bool `json:"manual"`
}

// OnColdStorage returns true if the disc is currently on cold storage (whatever storage grade that is)
func (s Snapshot) OnColdStorage() bool {
	return s.StorageGrade == ColdStorageGrade
}

// PrettyPrint outputs a nicely-formatted string detailing the snapshot to the given writer.
func (s Snapshot) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	snapshotTpl := `
{{ define "snapshot_sgl" }}{{ .Label }}{{ if not .OnColdStorage }} (in progress){{ end }}{{ end }}

{{ define "snapshot_medium" }}{{ template "snapshot_sgl" . }}{{ end }}

{{ define "snapshot_full" }}{{ template "snapshot_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, snapshotTpl, "snapshot"+string(detail), s)
}

// Snapshots represents a collection of snapshots
type Snapshots []*Snapshot

// PrettyPrint outputs a nicely-formatted string detailing the snapshot to the given writer.
func (ss Snapshots) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	snapshotsTpl := `
{{ define "snapshots_full" }}{{ template "snapshots_medium" . }}{{ end }}
{{ define "snapshots_medium" -}}
{{- range . -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "snapshots_sgl" }}{{ len . | pluralize "snapshot" "snapshots" }}{{ end }}
`
	return prettyprint.Run(wr, snapshotsTpl, "snapshots"+string(detail), ss)
}
