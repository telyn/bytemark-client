package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
)

// StoragePool represents a Bytemark Cloud Servers disk storage pool, as returned by the admin API.
type StoragePool struct {
	Name            string   `json:"name"`
	Label           string   `json:"label"`
	Zone            string   `json:"zone"`
	Size            int      `json:"size"`
	FreeSpace       int      `json:"free_space"`
	AllocatedSpace  int      `json:"alloc"`
	Discs           []string `json:"discs"`
	OvercommitRatio int      `json:"overcommit_ratio"`
	UsageStrategy   string   `json:"usage_strategy"`
	StorageGrade    string   `json:"grade"`
	Note            string   `json:"note"`

	// These were defined, but aren't returned by the API
	// ID        int
	// Tail      *Tail
	// IOPSLimit int
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
