package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
	"text/template"
)

const accountsTemplate = `{{ define "account_name" }}{{ if .BillingID }}{{ .BillingID }} - {{ end }}{{ if .Name }}{{ .Name }}{{ else }}[no bigv account]{{ end }}{{ end }}

{{ define "account_bullet" }}• {{ template "account_name" . -}}
{{- if .IsDefaultAccount }} (this is your default account){{ end -}}
{{- end }}

{{ define "whoami" }}You are '{{ .Username }}'{{ end }}

{{ define "owned_accounts" -}}
  {{- if .OwnedAccounts -}}
    Accounts you own: 
    {{- range .OwnedAccounts }}
  {{ template "account_bullet" . -}}
    {{ end -}}
  {{- end -}}
{{- end -}}

{{ define "extra_accounts" -}}
{{- if .OtherAccounts -}}
{{- if .OwnedAccounts }}

Other accounts you can access:
{{- else -}}
Accounts you can access:
{{- end -}}
{{- range .OtherAccounts }}
  {{template "account_bullet" . }}
{{- end -}}
{{- end -}}
{{- end }}

{{ define "group_overview" }}  • {{ .Name }} - {{  pluralize "server" "servers" ( len .VirtualMachines ) }}
{{ if ( len .VirtualMachines ) le 5 -}}
{{- range .VirtualMachines }}   {{ prettysprint . "_sgl" }}
{{ end -}}
{{- end -}}
{{- end -}}

{{/* account_overview needs $ to be defined, so use single_account_overview as entrypoint */}}
{{ define "account_full" }}
  {{- if .IsDefaultAccount -}}	
    Your default account ({{ template "account_name" . }})
  {{- else -}}
    {{- template "account_name" . -}}
  {{- end }}
{{ range .Groups -}}
    {{ template "group_overview" . -}}
{{- end -}}
{{- end }}

{{ define "single_account_overview" }}
{{ template "account_full" .Account }}
{{ end }}

{{ define "full_overview" -}}
{{- template "whoami" . }}

{{ template "owned_accounts" . -}}
{{- template "extra_accounts" . }}

{{ if .DefaultAccount -}}
{{- template "account_full" .DefaultAccount }}
{{ else -}}
It was not possible to determine your default account. Please set one using bytemark config set account.

{{ end -}}
{{- end }}
`

// FormatOverview outputs the given overview using the named template to the given writer.
// TODO(telyn): make template choice not a string
// TODO(telyn): use an actual Overview object?
func FormatOverview(wr io.Writer, accounts []*Account, defaultAccount *Account, username string) error {
	tmpl, err := template.New("accounts").Funcs(prettyprint.Funcs).Funcs(map[string]interface{}{
		"isDefaultAccount": func(a *Account) bool {
			if a == nil || defaultAccount == nil {
				return false
			}
			if a.BillingID != 0 && a.BillingID == defaultAccount.BillingID {
				return true
			}
			return a.Name != "" && a.Name == defaultAccount.Name
		},
	}).Parse(accountsTemplate)
	if err != nil {
		return err
	}
	var ownedAccounts []*Account
	var otherAccounts []*Account
	for _, a := range accounts {
		if a.Owner != nil && a.Owner.Username != "" && a.Owner.Username == username {
			ownedAccounts = append(ownedAccounts, a)
		} else {
			otherAccounts = append(otherAccounts, a)
		}
	}
	data := map[string]interface{}{
		"Accounts":       accounts,
		"DefaultAccount": defaultAccount,
		"Username":       username,
		"OwnedAccounts":  ownedAccounts,
		"OtherAccounts":  otherAccounts,
		"Writer":         wr,
	}

	err = tmpl.ExecuteTemplate(wr, "full_overview", data)
	return err
}
