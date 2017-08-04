package billing

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Person represents a bmbilling person
type Person struct {
	ID int `json:"id,omitempty"`
	// Username is the name this person uses to log in to our services.
	Username    string `json:"username"`
	Email       string `json:"email"`
	BackupEmail string `json:"email_backup,omitempty"`

	// only set in the creation request
	Password string `json:"password"`

	FirstName   string `json:"firstname"`
	LastName    string `json:"surname"`
	Address     string `json:"address"`
	City        string `json:"city"`
	StateCounty string `json:"statecounty,omitempty"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	MobilePhone string `json:"phonemobile,omitempty"`

	Organization         string `json:"organization,omitempty"`
	OrganizationDivision string `json:"division,omitempty"`
	VATNumber            string `json:"vatnumber,omitempty"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/bytemark-client for this type.
func (p Person) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "ID, Username, Email, FirstName, LastName, Phone, MobilePhone, Organization, OrganizationDivision"
	}
	return "ID, Username, Email, FirstName, LastName, Phone, MobilePhone, Address, City, Postcode, Country, Organization, OrganizationDivision, VATNumber"
}

// PrettyPrint ouputs the person to the writer in a human readable form, at the specified detail level.
func (p Person) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	personTpl := `
{{ define "person_sgl" }}{{ .FirstName }} {{ .LastName }} ({{ .Username }}) - {{ .Email }}{{ end }}
{{ define "person_medium" }}{{ template "person_sgl" . }}
Address: {{ .Address }}
    {{ .City }}
    {{ .Postcode }}
    {{ .Country }}

Phone: {{ .Phone }} 
Mobile: {{ .MobilePhone -}}
{{- end }}
{{ define "person_full" }}{{ template "person_medium" . }}

Organization: {{ .Organization }}
Division: {{ .OrganizationDivision }}
VAT Number: {{ .VATNumber -}}
{{- end }}
`
	return prettyprint.Run(wr, personTpl, "person"+string(detail), p)
}

// IsValid returns true if the Person is valid, false otherwise.
func (p Person) IsValid() bool {
	return p.Username != ""
}
