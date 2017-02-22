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

// Percent full gives us the (numeric) percentage of how full the disc is
func (sp StoragePool) PercentFull() int {
	return sp.AllocatedSpace * 100 / sp.Size
}

// PrettyPrint writes an overview of this account out to the given writer.
func (sp StoragePool) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const t = `
{{ define "storage_pool_sgl" -}}
• {{ .Label }}, Disc Count: {{ len .Discs }}, {{ .PercentFull }}% full
{{ end }}

{{ define "storage_pool_full" -}}
{{ template "storage_pool_sgl" . }}
{{- end }}`
	return prettyprint.Run(wr, t, "storage_pool"+string(detail), sp)
}
