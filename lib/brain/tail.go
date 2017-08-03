package brain

import (
	"io"
	"net"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Tail represents a Bytemark Cloud Servers tail (disk storage machine), as returned by the admin API.
type Tail struct {
	ID    int    `json:"id"`
	UUID  string `json:"uuid"`
	Label string `json:"label"`

	CCAddress net.IP `json:"cnc_address"`
	ZoneName  string `json:"zone"`

	IsOnline     bool     `json:"online"`
	StoragePools []string `json:"pools"`
}

func (t Tail) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "ID, Label, IsOnline, UUID, StoragePools, ZoneName"
	}
	return "ID, Label, IsOnline, UUID, CCAddress, StoragePools, ZoneName"
}

// PrettyPrint writes an overview of this tail out to the given writer.
func (t Tail) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const tpl = `
{{ define "tail_sgl" -}}
• {{ .Label }} (ID: {{ .ID }}), Online: {{ .IsOnline }}, Storage Pool Count: {{ len .StoragePools }}
{{ end }}

{{ define "tail_full" -}}
{{ template "tail_sgl" . }}
{{ template "storage_pools" . -}}
{{- end }}

{{ define "storage_pools"  }}
{{- if .StoragePools }}    storage pools:
{{- range .StoragePools }}
      • {{ . }}
{{- end }}

{{ end -}}
{{ end }}
`
	return prettyprint.Run(wr, tpl, "tail"+string(detail), t)
}
