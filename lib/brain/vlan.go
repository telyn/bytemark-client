package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
)

// VLAN is a representation of a VLAN, as used by admin endpoints
type VLAN struct {
	ID        int        `json:"id"`
	Num       int        `json:"num"`
	UsageType string     `json:"usage_type"`
	IPRanges  []*IPRange `json:"ip_ranges"`
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
