package lib

import (
	"bytes"
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Account represents both the BigV and bmbilling parts of an account.
type Account struct {
	Name             string         `json:"name"`
	Owner            billing.Person `json:"owner"`
	TechnicalContact billing.Person `json:"technical_contact"`
	BillingID        int            `json:"billing_id"`
	BrainID          int            `json:"brain_id"`
	CardReference    string         `json:"card_reference"`
	Groups           brain.Groups   `json:"groups"`
	Suspended        bool           `json:"suspended"`

	IsDefaultAccount bool `json:"-"`
}

func (a *Account) fillBrain(b brain.Account) {
	a.BrainID = b.ID
	a.Groups = b.Groups
	a.Suspended = b.Suspended
	a.Name = b.Name
}
func (a *Account) fillBilling(b billing.Account) {
	a.BillingID = b.ID
	a.Owner = b.Owner
	a.TechnicalContact = b.TechnicalContact
	a.CardReference = b.CardReference
	a.Name = b.Name
}

// CountVirtualMachines returns the number of virtual machines across all the Account's Groups.
func (a Account) CountVirtualMachines() (servers int) {
	for _, g := range a.Groups {
		servers += len(g.VirtualMachines)
	}
	return
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (a Account) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "BillingID, Name, Suspended"
	}
	return "BillingID, Name, Suspended, Groups"
}

// billingAccount copies all the billing parts of the account into a new billingAccount.
func (a Account) billingAccount() (b *billing.Account) {
	b = new(billing.Account)
	b.ID = a.BillingID
	b.Owner = a.Owner
	b.TechnicalContact = a.TechnicalContact
	b.CardReference = a.CardReference
	b.Name = a.Name
	return
}

// IsValid returns true if this account has been populated with data from the brain or billing endpoints.
func (a Account) IsValid() bool {
	return a.BillingID != 0 || a.BrainID != 0
}

// PrettyPrint writes an overview of this account out to the given writer.
func (a Account) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const accountsTemplate = `{{ define "account_name" }}{{ if .BillingID }}{{ .BillingID }} - {{ end }}{{ if .Name }}{{ .Name }}{{ else }}[no bigv account]{{ end }}{{ end }}

{{ define "account_sgl" }}• {{ template "account_name" . -}}
{{- if .IsDefaultAccount }} (this is your default account){{ end -}}
{{- end }}

{{ define "group_overview" }}  • {{ .Name }} - {{  pluralize "server" "servers" ( len .VirtualMachines ) -}}
{{- if len .VirtualMachines | gt 6 -}}
{{- range .VirtualMachines }}
   {{ prettysprint . "_sgl" -}}
{{- end -}}
{{- end }}
{{ end -}}

{{ define "account_medium" }}
  {{- if .IsDefaultAccount -}}	
    Your default account ({{ template "account_name" . }})
  {{- else -}}
    {{- template "account_name" . -}}
  {{- end -}}
{{ end }}

{{ define "account_full" -}}
{{- template "account_medium" . }}
{{ range .Groups -}}
    {{ template "group_overview" . -}}
{{- end -}}
{{- end }}`
	return prettyprint.Run(wr, accountsTemplate, "account"+string(detail), a)
}

// String formats this account as a string - the same format as prettyprint.SingleLine
func (a Account) String() string {
	buf := bytes.Buffer{}
	err := a.PrettyPrint(&buf, prettyprint.SingleLine)
	if err != nil {
		return "ERROR"
	}
	return buf.String()
}

type Accounts []Account

func (as Accounts) DefaultFields(f output.Format) string {
	return (Account{}).DefaultFields(f)
}

func (as Accounts) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	accountsTpl := `
{{ define "accounts_sgl" -}}
{{- range . -}}
{{- .Name }}, {{ end -}}
{{- end }}

{{ define "accounts_medium" }}{{ template "accounts_sgl" . }}{{ end }}

{{ define "accounts_full" }}
Accounts: 
{{ range . -}}
{{- prettysprint "_sgl" . }}
{{ end -}}
{{- end }}`
	return prettyprint.Run(wr, accountsTpl, "accounts"+string(detail), as)
}
