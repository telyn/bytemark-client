package billing

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Definition is an admin-modifiable parameter for bmbilling
// examples include account-opening credit amount and trial length
type Definition struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Value string `json:"value"`
	// Which auth group a user must be in to update the definition
	UpdateGroupReq string `json:"update_group_req,omitempty"`
}

// DefaultFields returns the default fields used for making tables of Definitions
func (d Definition) DefaultFields(f output.Format) string {
	return "Name, Value, UpdateGroupReq"
}

// PrettyPrint writes the Definition in a human-readable form at the given detail level to wr
func (d Definition) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) {
	definitionTpl := `
		{{ define "definition_sgl" }}{{ .Name }}: {{ .Value }}{{ end }}
		{{ define "definition_medium" }}{{ template "definition_sgl" . }}{{ end }}
		{{ define "definition_full" -}}
ID: {{ .ID }}
Name: {{ .Name }}
Value: {{ .Value }}
Update Group Requirement: {{ .UpdateGroupReq }}
		{{- end }}
	`
	prettyprint.Run(wr, definitionTpl, "definition"+detail)
}
