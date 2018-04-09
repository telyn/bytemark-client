package brain

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// User represents a Bytemark user.
type User struct {
	Username       string
	Email          string
	AuthorizedKeys Keys
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (user User) DefaultFields(f output.Format) string {
	return "Username, Email"
}

// PrettyPrint outputs human-readable information about the user to the given writer at some level of detail.
func (user User) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	userTpl := `
{{ define "user_sgl" }}{{ .Username }}{{ end }}
{{ define "user_medium" }}{{ .Username }} - {{ .Email }}{{ end }}
{{ define "user_full " }}{{ template "user_medium" }}

Authorized keys:
{{ range .AuthorizedKeys }}
{{ . }}	
{{ end }}
`
	return prettyprint.Run(wr, userTpl, "user"+string(detail), user)
}
