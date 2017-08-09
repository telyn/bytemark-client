package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// VirtualMachines represents more than one account in output.Outputtable form.
type VirtualMachines []VirtualMachine

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type, which is the same as VirtualMachine.DefaultFields.
func (gs VirtualMachines) DefaultFields(f output.Format) string {
	return (VirtualMachine{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the virtual machines to writer at the given detail level.
func (gs VirtualMachines) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	virtualmachinesTpl := `
{{ define "virtualmachines_sgl" }}{{ len . }} servers{{ end }}

{{ define "virtualmachines_medium" -}}
{{- range . -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "virtualmachines_full" }}{{ template "virtualmachines_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, virtualmachinesTpl, "virtualmachines"+string(detail), gs)
}
