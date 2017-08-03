package brain

// Account represents an account object that's returned by the brain
type Account struct {
	Name string `json:"name"`

	// the following cannot be set
	ID        int     `json:"id"`
	Suspended bool    `json:"suspended"`
	Groups    []Group `json:"groups"`
}

func (a Account) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "Name", "Suspended"
	}
	return "ID", "Name", "Groups"
}

func (a Account) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) {
	accountTpl := `
	{{ define "account_sgl" }}{{ .Name }}{{ if .Suspended }} (suspended){{ end}}{{ end }}
{{ define "account_medium" }}{{ template "account_sgl" . }}{{ end }}
{{ define "account_full" }}{{ .Name }}{{ if .Suspended }} (suspended){{ end}}{{ end }}

{{ range .Groups -}}
    {{ template "group_overview" . -}}
{{- end -}}
{{ define "group_overview" }}  • {{ .Name }} - {{  pluralize "server" "servers" ( len .VirtualMachines ) -}}
{{- if len .VirtualMachines | gt 6 -}}
{{- range .VirtualMachines }}
   {{ prettysprint . "_sgl" -}}
{{- end -}}
{{- end }}
{{ end -}}
	`
	return prettyprint.Run(wr, accountTpl, "account"+string(detail))
}
