package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// TODO(tom): add test

// VMDefault represents a VM Default, as passed around from the vm_defaults endpoint
type VMDefault struct {
	CdromURL        string `json:"cdrom_url,omitempty"`
	Cores           int    `json:"cores,omitempty"`
	Memory          int    `json:"memory,omitempty"`
	Name            string `json:"name,omitempty"`
	HardwareProfile string `json:"hardware_profile,omitempty"`

	// zone name can be set during creation but not changed
	ZoneName string `json:"zone_name,omitempty"`

	// the following cannot be set
	Discs Discs `json:"discs,omitempty"`
	ID    int   `json:"id,omitempty"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (vmd VMDefault) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "Memory, Cores, Name "
	}
	return "Memory, Cores, Discs, CdromURL, HardwareProfile"
}

// PrettyPrint outputs a nice human-readable overview of the VM Default to the given writer.
func (vmd VMDefault) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = ` 
{{ define "vmdefault_sgl" }} ▸ {{.Name }} in {{capitalize .ZoneName}}{{ end }}
{{ define "vmdefault_spec" }}   - {{ pluralize "core" "cores" .Cores }}, {{ mibgib .Memory }}, {{ if .Discs}}{{.TotalDiscSize "" | gibtib }} on {{ len .Discs | pluralize "disc" "discs"  }}{{ else }}no discs{{ end }}{{ end }}

{{ define "vmdefault_discs"  }}
{{- if .Discs }}    discs:
{{- range .Discs }}
      • {{ prettysprint . "_sgl" }}
{{- end }}

{{ end -}}
{{ end }}

{{ define "vmdefault_medium" }}{{ template "vmdefault_sgl" . }}
{{ template "vmdefault_spec" . }}{{ end }}

{{ define "vmdefault_full" -}}
{{ template "vmdefault_medium" . }}

{{ template "vmdefault_discs" . -}}
{{ end }}`
	return prettyprint.Run(wr, template, "vmdefault"+string(detail), vmd)
}

// TotalDiscSize returns the sum of all disc capacities in the VM for the given storage grade.
// Provide the empty string to sum all discs regardless of storage grade.
func (vmd VMDefault) TotalDiscSize(storageGrade string) (total int) {
	total = 0
	for _, disc := range vmd.Discs {
		if storageGrade == "" || storageGrade == disc.StorageGrade {
			total += disc.Size
		}
	}
	return total
}
