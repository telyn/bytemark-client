package lib

import (
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
	"text/template"
)

// FormatOverview outputs the given overview using the named template to the given writer.
// TODO(telyn): make template choice not a string
// TODO(telyn): use an actual Overview object?
func FormatOverview(wr io.Writer, accounts []*Account, defaultAccount *Account, username string) error {
	const overviewTemplate = `{{ define "account_name" }}{{ if .BillingID }}{{ .BillingID }} - {{ end }}{{ if .Name }}{{ .Name }}{{ else }}[no bigv account]{{ end }}{{ end }}

{{ define "whoami" }}You are '{{ .Username }}'{{ end }}

{{ define "owned_accounts" -}}
  {{- if .OwnedAccounts -}}
    Accounts you own: 
    {{- range .OwnedAccounts }}
  {{ prettysprint . "_sgl" -}}
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
  {{ prettysprint . "_sgl" }}
{{- end -}}
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
{{- prettysprint .DefaultAccount "_full" }}
{{ else -}}
It was not possible to determine your default account. Please set one using bytemark config set account.

{{ end -}}
{{- end }}
`
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
	}).Parse(overviewTemplate)
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
