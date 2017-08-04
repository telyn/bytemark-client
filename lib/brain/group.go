package brain

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"io"
)

// Group represents a group
type Group struct {
	Name string `json:"name"`

	// the following cannot be set
	AccountID       int              `json:"account_id"`
	ID              int              `json:"id"`
	VirtualMachines []VirtualMachine `json:"virtual_machines"`
}

func (g Group) CountVirtualMachines() int {
	return len(g.VirtualMachines)
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/bytemark-client for this type.
func (g Group) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "Name, CountVirtualMachines"
	}
	return "Name, VirtualMachines"
}

func (g Group) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	groupTpl := `
{{ define "group_sgl " -}}
  {{- .Name }} - Group containing {{ len .VirtualMachines }} cloud {{ len .VirtualMachines | pluralize "server" "servers" -}}
{{ end }}

{{ define "group_medium" -}}
{{- template "group_sgl" . -}}
{{- end }}

{{ define "group_full" -}}
{{ template "group_sgl" . }}

{{ for .VirtualMachines -}}
{{- prettysprint "_sgl" . }}
{{ end -}}
{{- end }}
`
	return prettyprint.Run(wr, groupTpl, "group"+string(detail), g)
}

func (g Group) String() string {
	return fmt.Sprintf("group %d %q - has %d servers", g.ID, g.Name, len(g.VirtualMachines))
}
