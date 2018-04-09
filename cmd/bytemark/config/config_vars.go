package config

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

var configVars = [...]string{
	"endpoint",
	"billing-endpoint",
	"auth-endpoint",
	"spp-endpoint",
	"admin",
	"user",
	"account",
	"group",
	"output-format",
	"session-validity",
	"token",
	"debug-level",
	"yubikey",
}

// IsConfigVar checks to see if the named variable is actually one of the settable configVars.
func IsConfigVar(name string) bool {
	for _, v := range configVars {
		if v == name {
			return true
		}
	}
	return false
}

// Vars is a list of configuration variables
type Vars []Var

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type, which is the same as VirtualMachine.DefaultFields.
func (vs Vars) DefaultFields(f output.Format) string {
	return (Var{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the virtual machines to writer at the given detail level.
func (vs Vars) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{ define "vars_sgl" }}{{ len . }} variables{{end}}{{define "vars_medium" }}{{ range . }}{{ prettysprint . "_sgl" }}{{ end }}{{ define "vars_full"}}{{ template "vars_medium" . }}{{ end }}`
	return prettyprint.Run(wr, template, "vars"+string(detail), vs)
}
