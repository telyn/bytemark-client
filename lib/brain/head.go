package brain

import (
	"io"
	"net"

	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
)

// Head represents a Bytemark Cloud Servers head server.
type Head struct {
	ID       int    `json:"id,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	Label    string `json:"label,omitempty"`
	ZoneName string `json:"zone,omit_empty"`

	// descriptive information that can be modified

	Architecture  string   `json:"arch"`
	CCAddress     *net.IP  `json:"cnc_address"`
	Note          string   `json:"note"`
	Memory        int      `json:"memory,omitempty"`
	UsageStrategy string   `json:"usage_strategy,omitempty"`
	Models        []string `json:"models,omitempty"`

	// state

	MemoryFree int  `json:"memory_free,omitempty"`
	IsOnline   bool `json:"is_online,omitempty"`
	UsedCores  int  `json:"used_cores"`

	// You may have one or the other.

	VirtualMachineCount int      `json:"virtual_machines_count,omitempty"`
	VirtualMachines     []string `json:"virtual_machines,omitempty"`
}

// CountVirtualMachines returns the number of virtual machines running on this head
func (h *Head) CountVirtualMachines() int {
	if h.VirtualMachines != nil {
		return len(h.VirtualMachines)
	}
	return h.VirtualMachineCount
}

// PrettyPrint writes an overview of this head out to the given writer.
func (h Head) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const t = `
{{ define "head_sgl" -}}
• {{ .Label }} (ID: {{ .ID }}), Online: {{ .IsOnline }}, VM Count: {{ len .VirtualMachines }}
{{ end }}

{{ define "head_full" -}}
{{ template "head_sgl" . }}
{{- end }}`
	return prettyprint.Run(wr, t, "head"+string(detail), h)
}
