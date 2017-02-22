package brain

import (
	"io"
	"net"

	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
)

// Tail represents a Bytemark Cloud Servers tail (disk storage machine), as returned by the admin API.
type Tail struct {
	ID    int    `json:"id"`
	UUID  string `json:"uuid"`
	Label string `json:"label"`

	CCAddress *net.IP `json:"cnc_address"`
	ZoneName  string  `json:"zone"`

	IsOnline     bool     `json:"online"`
	StoragePools []string `json:"pools"`
}

// PrettyPrint writes an overview of this tail out to the given writer.
func (t Tail) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const tpl = `
{{ define "tail_sgl" -}}
• {{ .Label }} (ID: {{ .ID }}), Online: {{ .IsOnline }}, Storage Pool Count: {{ len .StoragePools }}
{{ end }}

{{ define "tail_full" -}}
{{ template "tail_sgl" . }}
{{- end }}`
	return prettyprint.Run(wr, tpl, "tail"+string(detail), t)
}
