package brain

import (
	"io"
	"net"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// IPCreateRequest is used by the create_ip endpoint on a nic
type IPCreateRequest struct {
	Addresses  int    `json:"addresses"`
	Family     string `json:"family"`
	Reason     string `json:"reason"`
	Contiguous bool   `json:"contiguous"`
	// don't actually specify the IPs - this is for filling in from the response!
	IPs []net.IP `json:"ips"`
}

// DefaultFields returns the default fields for feeding into github.com/BytemarkHosting/row.FieldsFrom
func (ipcr IPCreateRequest) DefaultFields(f output.Format) string {
	return "Addresses, Family, Reason, Contiguous, IPs"
}

// PrettyPrint outputs this create request/response in a human-readable form
// TODO: support the response
func (ipcr IPCreateRequest) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	ipcrTpl := `
{{ define "ipcr_sgl" }}{{ .Addresses }} {{ .Family }} {{ if .Contiguous }}contiguous {{ end }}addresses wanted for {{ .Reason }}{{ end }}
{{ define "ipcr_medium" }}{{ template "ipcr_sgl" . }}{{ end }}
{{ define "ipcr_full" }}{{ template "ipcr_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, ipcrTpl, "ipcr"+string(detail), ipcr)
}
