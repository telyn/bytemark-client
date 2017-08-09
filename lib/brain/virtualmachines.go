package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

type VirtualMachines []VirtualMachine

func (gs VirtualMachines) DefaultFields(f output.Format) string {
	return (VirtualMachine{}).DefaultFields(f)
}

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
