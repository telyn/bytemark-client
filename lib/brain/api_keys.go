package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// APIKeys is a slice of APIKey implementing output.Outputtable
type APIKeys []APIKey

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (keys APIKeys) DefaultFields(f output.Format) string {
	return (APIKey{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the API keys to the given
// writer.
func (keys APIKeys) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	keysTpl := `
{{ define "keys_sgl" }}{{ len . }} API keys{{ end }}
{{ define "keys_medium" }}
{{- range . -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{ end }}
{{ define "keys_full" }}{{ template "keys_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, keysTpl, "keys"+string(detail), keys)
}
