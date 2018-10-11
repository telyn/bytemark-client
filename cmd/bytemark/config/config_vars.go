package config

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

var configVars = [...]string{
	"account",
	"admin",
	"auth-endpoint",
	"billing-endpoint",
	"debug-level",
	"endpoint",
	"group",
	"insecure",
	"output-format",
	"session-validity",
	"spp-endpoint",
	"token",
	"user",
	"yubikey",
}

// VarsDescription is text suitable for inclusion in commands that manipulate
// config variables.
const VarsDescription = `
        account - the default account, used when you do not explicitly state an account - defaults to the same as your user name
        token - the token used for authentication
        user - the user that you log in as by default
        group - the default group, used when you do not explicitly state a group (defaults to 'default')

        debug-level - the default debug level. Set to 0 unless you like lots of output.
	api-endpoint - the endpoint for domains (among other things?)
        auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.
        endpoint - the brain endpoint to connect to. https://uk0.bigv.io is the default.
        billing-endpoint - the billing API endpoint to connect to. https://bmbilling.bytemark.co.uk is the default.
        spp-endpoint - the SPP endpoint to use. https://spp-submissions.bytemark.co.uk is the default.`

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
