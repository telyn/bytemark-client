package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// StoragePool represents a Bytemark Cloud Servers disk storage pool, as returned by the admin API.
type StoragePool struct {
	Name                 string `json:"name,omitempty"`
	Label                string `json:"label,omitempty"`
	ZoneName             string `json:"zone,omitempty"`
	Size                 int    `json:"size,omitempty"`
	FreeSpace            int    `json:"free_space,omitempty"`
	Ceiling              int    `json:"ceiling,omitempty"`
	AllocatedSpace       int    `json:"alloc,omitempty"`
	Discs                int    `json:"discs,omitempty"`
	Backups              int    `json:"backups,omitempty"`
	OvercommitRatio      int    `json:"overcommit_ratio,omitempty"`
	MigrationConcurrency int    `json:"migration_concurrency,omitempty"`
	UsageStrategy        string `json:"usage_strategy,omitempty"`
	StorageGrade         string `json:"grade,omitempty"`
	Note                 string `json:"note,omitempty"`
	IOPSLimit            int    `json:"iops_limit,omitempty"`

	// These were defined, but aren't returned by the API
	// ID        int
	// Tail      Tail
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (sp StoragePool) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "ZoneName, Label, Ceiling, Size, FreeSpace, AllocatedSpace, Discs, OvercommitRatio, UsageStrategy, MigrationConcurrency, Note"
	}
	return "ZoneName, Label, Ceiling, Size, FreeSpace, AllocatedSpace, Discs, OvercommitRatio, UsageStrategy, MigrationConcurrency, Note"
}

// PercentFull gives us the (numeric) percentage of how full the disc is
func (sp StoragePool) PercentFull() int {
	// If no space is allocated, the disc is full, so return 100
	if sp.Size == 0 {
		return 100
	}

	return sp.AllocatedSpace * 100 / sp.Size
}

// PrettyPrint writes an overview of this storage pool out to the given writer.
func (sp StoragePool) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const t = `
{{ define "storage_pool_sgl" -}}
• {{ .Label }}, Disc Count: {{ len .Discs }}, {{ .PercentFull }}% full
{{ end }}

{{ define "storage_pool_full" -}}
{{ template "storage_pool_sgl" . }}
{{ template "discs" . }}
{{- end }}

{{ define "discs"  }}
{{- if .Discs }}    discs:
{{- range .Discs }}
      • {{ . }}
{{- end }}

{{ end -}}
{{ end }}
`
	return prettyprint.Run(wr, t, "storage_pool"+string(detail), sp)
}

// StoragePools represents more than one storage pool in output.Outputtable form.
type StoragePools []StoragePool

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type, which is the same as StoragePool.DefaultFields
func (sps StoragePools) DefaultFields(f output.Format) string {
	return (StoragePool{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the storage pools to writer at the given detail level.
func (sps StoragePools) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	storagepoolsTpl := `
{{ define "storagepools_sgl" }}{{ len . }} servers{{ end }}

{{ define "storagepools_medium" -}}
{{- range . -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "storagepools_full" }}{{ template "storagepools_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, storagepoolsTpl, "storagepools"+string(detail), sps)
}
