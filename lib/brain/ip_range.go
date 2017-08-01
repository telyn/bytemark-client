package brain

import (
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// IPRange is a representation of an IP range
type IPRange struct {
	ID        int      `json:"id"`
	Spec      string   `json:"spec"`
	VLANNum   int      `json:"vlan_num"`
	Zones     []string `json:"zones"`
	Available *big.Int `json:"available"` // wants to be a pointer because MarshalText is defined on the pointer type and we need it for the tests (but not for non-tests)
}

// String serialises an IP range to easily be output
func (ipr IPRange) String() string {
	zones := ""
	if len(ipr.Zones) > 0 {
		pluralise := ""
		if len(ipr.Zones) > 1 {
			pluralise = "s"
		}
		zones = fmt.Sprintf(", in zone%s %s", pluralise, strings.Join(ipr.Zones, ","))
	}
	// Since `Available` is a float64 but won't have decimal points, we just format accordingly.
	return fmt.Sprintf("%s%s, %.0f IPs available.", ipr.Spec, zones, ipr.Available)
}

// PrettyPrint writes an overview of this IP range out to the given writer.
func (ipr IPRange) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const t = `
{{ define "ip_range_sgl" -}}
• IP range {{ .Spec }} (ID: {{ .ID }}), has {{ len .Zones }} zones and {{ printf "%.0f" .Available }} IPs available.
{{ end }}

{{ define "ip_range_full" -}}
{{ template "ip_range_sgl" . }}
{{ template "zones" . }}
{{- end }}

{{ define "zones"  }}
{{- if .Zones }}    zones:
{{- range .Zones }}
      • {{ . }}
{{- end }}

{{ end -}}
{{ end }}
`
	return prettyprint.Run(wr, t, "ip_range"+string(detail), ipr)
}
