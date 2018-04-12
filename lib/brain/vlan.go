package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// VLAN is a representation of a VLAN, as used by admin endpoints
type VLAN struct {
	ID        int      `json:"id"`
	Num       int      `json:"num"`
	UsageType string   `json:"usage_type"`
	IPRanges  IPRanges `json:"ip_ranges"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (v VLAN) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "ID, Num, UsageType"
	}
	return "ID, Num, UsageType, IPRanges"
}

// PrettyPrint writes an overview of this VLAN out to the given writer.
func (v VLAN) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const t = `
{{ define "vlan_sgl" -}}
• {{ .UsageType }} (Num: {{ .Num }}). Contains {{ len .IPRanges }} IP ranges.
{{ end }}

{{ define "vlan_full" -}}
{{ template "vlan_sgl" . }}
{{ template "ip_ranges" . }}
{{- end }}

{{ define "ip_ranges"  }}
{{- if .IPRanges }}    IP ranges:
{{- range .IPRanges }}
      {{ template "ip_range" . }}
{{- end }}

{{ end -}}
{{ end }}

{{ define "ip_range" -}}
• {{ .String }}
{{- end }}
`
	return prettyprint.Run(wr, t, "vlan"+string(detail), v)
}

// VLANs represents more than one VLAN in output.Outputtable form.
type VLANs []VLAN

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type, which is the same as VLAN.DefaultFields
func (vs VLANs) DefaultFields(f output.Format) string {
	return (VLAN{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the VLANs to writer at the given detail level.
func (vs VLANs) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	vlansTpl := `
{{ define "vlans_sgl" }}{{ len . }} servers{{ end }}

{{ define "vlans_medium" -}}
{{- range -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "vlans_full" }}{{ template "vlans_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, vlansTpl, "vlans"+string(detail), vs)
}
