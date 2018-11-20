package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// VMDefaultSpec represents a VM Default specification.
type VMDefaultSpec struct {
	VMDefault VMDefault     `json:"vm_default,omitempty"`
	Discs     []Disc        `json:"disc,omitempty"`
	Reimage   *ImageInstall `json:"reimage,omitempty"` // may want to be null, so is a pointer
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (spec VMDefaultSpec) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "VirtualMachine, Discs, Reimage"
	}
	return "VirtualMachine, Discs, Reimage"
}

// PrettyPrint outputs a nice human-readable overview of the VM Default Specification to the given writer.
func (spec VMDefaultSpec) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{ define "vmdspec_sgl" }} ▸ {{.VMDefault.Name }} in {{capitalize .VMDefault.ZoneName}}{{ end }}
{{ define "vmdspec_vmd" }} {{ pluralize "core" "cores" .VMDefault.Cores }}, {{ mibgib .VMDefault.Memory }}, {{pluralize "disc" "discs"  }}{{ else }}no discs{{ end }}{{ end }}
{{ define "vmdspec_reimage" }} {{ .Reimage.Distribution }}{{ end }}

{{ define "vmdspec_discs"  }}
{{- if .Discs }}    discs:
{{- range .VMDefault.Discs }}
      • {{ prettysprint . "_sgl" }}
{{- end }}
{{ end -}}
{{ end }}

{{ define "vmdspec_medium" }}{{ template "vmdspec_sgl" . }}
{{ template "vmdspec_vmd" . }}{{ end }}

{{ define "vmdspec_full" -}}
{{ template "vmdspec_medium" . }}

{{ template "vmdspec_vmd" . }}
{{ template "vmdspec_discs" . -}}
{{ template "vmdspec_reimage" . }}
{{ end }}`
	return prettyprint.Run(wr, template, "vmdspec"+string(detail), spec)
}
