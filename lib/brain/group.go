package brain

import (
	"fmt"
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Group represents a group
type Group struct {
	Name string `json:"name"`

	// the following cannot be set
	AccountID       int              `json:"account_id"`
	ID              int              `json:"id"`
	VirtualMachines []VirtualMachine `json:"virtual_machines"`
}

// CountVirtualMachines returns the number of virtual machines in this group
func (g Group) CountVirtualMachines() int {
	return len(g.VirtualMachines)
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (g Group) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "Name, CountVirtualMachines"
	}
	return "Name, VirtualMachines"
}

// PrettyPrint outputs a vaguely human-readable version of the definition to wr.
func (g Group) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	groupTpl := `
{{ define "group_sgl" -}}
  {{- .Name }} - Group containing {{ len .VirtualMachines }} cloud {{ len .VirtualMachines | pluralize "server" "servers" -}}
{{ end }}

{{ define "group_medium" -}}
{{- template "group_sgl" . -}}
{{- end }}

{{ define "group_full" -}}
{{ template "group_sgl" . }}

{{ range .VirtualMachines -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}
`
	return prettyprint.Run(wr, groupTpl, "group"+string(detail), g)
}

// String formats the Group as a string - a single line in a human-readable form.
func (g Group) String() string {
	return fmt.Sprintf("group %d %q - has %d servers", g.ID, g.Name, len(g.VirtualMachines))
}

// Groups represents multiple Group objects in output.Outputtable form.
type Groups []Group

// DefaultFields provides the default fields for Groups, which is the same as those of Group.
func (gs Groups) DefaultFields(f output.Format) string {
	return (Group{}).DefaultFields(f)
}

// PrettyPrint writes a summary of the groups to wr at the provided detail level.
func (gs Groups) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	groupsTpl := `
{{ define "groups_sgl" }}{{ len . }} groups{{ end }}

{{ define "groups_medium" -}}
{{- range -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "groups_full" }}{{ template "groups_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, groupsTpl, "groups"+string(detail), gs)
}

func (g Group) Path() (string, error) {
	return GroupID(g.ID).Path()
}
