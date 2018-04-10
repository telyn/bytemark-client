package config

import (
	"io"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Var is a struct which contains a name-value-source triplet
// Source is up to two words separated by a space. The first word is the source type: FLAG, ENV, DIR, CODE.
// The second is the name of the flag/file/environment var used.
type Var struct {
	Name   string
	Value  string
	Source string
}

// SourceType returns one of the following:
// FLAG for a configVar whose value was set by passing a flag on the command line
// ENV for a configVar whose value was set from an environment variable
// DIR for a configVar whose value was set from a file in the config dir
//
func (v *Var) SourceType() string {
	bits := strings.Fields(v.Source)

	return bits[0]
}

// SourceBaseName returns the basename of the configVar's source.
// it's a bit stupid and so its output is only valid for configVars with SourceType() of DIR
func (v *Var) SourceBaseName() string {
	bits := strings.Split(v.Source, "/")
	return bits[len(bits)-1]
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (v Var) DefaultFields(f output.Format) string {
	return "Name, Value, Source"
}

// PrettyPrint outputs a nice human-readable overview of the server to the given writer.
func (v Var) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{ define "var_sgl" }} ▸ {{.Name}} {{.Value}} ({{.Source}}){{ end }}{{ define "var_medium" }}{{ template "var_sgl" . }}{{ end }}{{ define "var_full" }}{{ template "var_sgl" . }}{{ end }}`
	return prettyprint.Run(wr, template, "var"+string(detail), v)
}
