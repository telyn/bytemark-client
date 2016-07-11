package lib

import (
	"fmt"
	"io"
	"math"
	"strings"
	"text/template"
	"unicode"
)

const accountsTemplate = `{{ define "account_name" }}{{ if .BillingID }}{{ .BillingID }} - {{ end }}{{ if .Name }}{{ .Name }}{{ else }}[no bigv account]{{ end }}{{ end }}

{{ define "account_bullet" }}• {{ template "account_name" . -}}
{{- if isDefaultAccount . }} (this is your default account){{ end -}}
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
{{- else }}
Accounts you can access:
{{- end -}}
{{- range .OtherAccounts }}
  {{template "account_bullet" . }}
{{- end -}}
{{- end -}}
{{- end }}

{{ define "group_overview" }}  • {{ .Name }} - {{  pluralize "server" "servers" ( len .VirtualMachines ) }}
{{ if ( len .VirtualMachines ) le 5 -}}
{{- range .VirtualMachines }}   {{ template "server_summary" . }}
{{ end -}}
{{- end -}}
{{- end -}}

{{/* account_overview needs $ to be defined, so use single_account_overview as entrypoint */}}
{{ define "account_overview" }}
  {{- if isDefaultAccount . -}}	
    Your default account ({{ template "account_name" . }})
  {{- else -}}
    {{- template "account_name" . -}}
  {{- end }}
{{ range .Groups -}}
    {{ template "group_overview" . -}}
{{- end -}}
{{- end }}

{{ define "single_account_overview" }}
{{ template "account_overview" .Account }}
{{ end }}

{{ define "full_overview" -}}
{{- template "whoami" . }}

{{ template "owned_accounts" . }}
{{ template "extra_accounts" . }}

{{ template "account_overview" .DefaultAccount }}
{{ end }}
`

const serverTemplate = `{{ define "server_summary" }} ▸ {{.ShortName }} ({{ if .Deleted }}deleted{{ else if .PowerOn }}powered on{{else}}powered off{{end}}) in {{capitalize .ZoneName}}{{ end }}
{{ define "server_spec" }}   {{ .PrimaryIP }} - {{ pluralize "core" "cores" .Cores }}, {{ mibgib .Memory }}, {{.TotalDiscSize "" | gibtib }} on {{ len .Discs | pluralize "disc" "discs"  }}{{ end }}

{{ define "discs" }}    discs:
{{- range .Discs }}
      • {{ .Label }} - {{ gibtib .Size }}, {{ .StorageGrade }} grade
{{- end }}
{{ end }}

{{ define "ips" }}    ips:
{{- range .AllIPv4Addresses }}
      • {{.}}
{{- end}}
{{- range .AllIPv6Addresses }}
      • {{.}}
{{- end }}
{{ end }}

{{ define "server_twoline" }}{{ template "server_summary" . }}
{{ template "server_spec" . }}{{ end }}

{{ define "server_full" -}}
{{ template "server_twoline" . }}

{{ template "discs" . }}
{{ template "ips" . }}
{{ end }}`

type TemplateChoice string

const (
	OneLine TemplateChoice = "server_summary"
	TwoLine                = "server_twoline"
	All                    = "server_full"
)

var templateFuncMap = map[string]interface{}{
	"mibgib": func(size int) string {
		mg := 'M'
		if size >= 1024 {
			size /= 1024
			mg = 'G'
		}
		return fmt.Sprintf("%d%ciB", size, mg)
	},
	"gibtib": func(size int) string {
		// lt is a less than sign if < 1GiB
		lt := ""
		if size < 1024 {
			lt = "< "
		}
		size /= 1024
		gt := 'G'
		if size >= 1024 {
			size /= 1024
			gt = 'T'
		}
		return fmt.Sprintf("%s%d%ciB", lt, size, gt)
	},
	"capitalize": func(str string) string {
		if len(str) == 0 {
			return str
		}

		runes := []rune(str)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	},
	"pluralize": func(single string, plural string, num int) string {
		if num == 1 {
			return fmt.Sprintf("%d %s", num, single)
		}
		return fmt.Sprintf("%d %s", num, plural)
	},
}

func FormatVirtualMachine(wr io.Writer, vm *VirtualMachine, tpl TemplateChoice) error {
	tmpl, err := template.New("virtualmachine").Funcs(templateFuncMap).Parse(serverTemplate)
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(wr, string(tpl), vm)
	if err != nil {
		return err
	}

	return nil
}

func FormatImageInstall(wr io.Writer, ii *ImageInstall, tpl TemplateChoice) error {
	output := make([]string, 0)
	if ii.Distribution != "" {
		output = append(output, "Image: "+ii.Distribution)
	}
	if ii.PublicKeys != "" {
		keynames := make([]string, 0)
		for _, k := range strings.Split(ii.PublicKeys, "\n") {
			kbits := strings.SplitN(k, " ", 3)
			if len(kbits) == 3 {
				keynames = append(keynames, strings.TrimSpace(kbits[2]))
			}

		}
		output = append(output, fmt.Sprintf("%d public keys: %s", len(keynames), strings.Join(keynames, ", ")))
	}
	if ii.RootPassword != "" {
		output = append(output, "Root/Administrator password: "+ii.RootPassword)
	}
	if ii.FirstbootScript != "" {
		output = append(output, "With a firstboot script")
	}
	_, err := wr.Write([]byte(strings.Join(output, "\r\n")))
	return err
}

func FormatVirtualMachineSpec(wr io.Writer, group *GroupName, spec *VirtualMachineSpec, tpl TemplateChoice) error {
	output := make([]string, 0, 10)
	output = append(output, fmt.Sprintf("Name: '%s'", spec.VirtualMachine.Name))
	output = append(output, fmt.Sprintf("Group: '%s'", group.Group))
	if group.Account == "" {
		output = append(output, "Account: not specified - will default to the account with the same name as the user you log in as")
	} else {
		output = append(output, fmt.Sprintf("Account: '%s'", group.Account))
	}
	s := ""
	if spec.VirtualMachine.Cores > 1 {
		s = "s"
	}

	mems := fmt.Sprintf("%d", spec.VirtualMachine.Memory/1024)
	if 0 != math.Mod(float64(spec.VirtualMachine.Memory), 1024) {
		mem := float64(spec.VirtualMachine.Memory) / 1024.0
		mems = fmt.Sprintf("%.2f", mem)
	}
	output = append(output, fmt.Sprintf("Specs: %d core%s and %sGiB memory", spec.VirtualMachine.Cores, s, mems))

	locked := ""
	if spec.VirtualMachine.HardwareProfile != "" {
		if spec.VirtualMachine.HardwareProfileLocked {
			locked = " (locked)"
		}
		output = append(output, fmt.Sprintf("Hardware profile: %s%s", spec.VirtualMachine.HardwareProfile, locked))
	}

	if spec.IPs != nil {
		if spec.IPs.IPv4 != "" {
			output = append(output, fmt.Sprintf("IPv4 address: %s", spec.IPs.IPv4))
		}
		if spec.IPs.IPv6 != "" {
			output = append(output, fmt.Sprintf("IPv6 address: %s", spec.IPs.IPv6))
		}
	}

	if spec.Reimage != nil {
		if spec.Reimage.Distribution == "" {
			if spec.VirtualMachine.CdromURL == "" {
				output = append(output, "No image or CD URL specified")
			} else {
				output = append(output, fmt.Sprintf("CD URL: %s", spec.VirtualMachine.CdromURL))
			}
		} else {
			output = append(output, "Image: "+spec.Reimage.Distribution)
		}
		output = append(output, "Root/Administrator password: "+spec.Reimage.RootPassword)
	} else {

		if spec.VirtualMachine.CdromURL == "" {
			output = append(output, "No image or CD URL specified")
		} else {
			output = append(output, fmt.Sprintf("CD URL: %s", spec.VirtualMachine.CdromURL))
		}
	}

	s = ""
	if len(spec.Discs) > 1 {
		s = "s"
	}
	if len(spec.Discs) > 0 {
		output = append(output, fmt.Sprintf("%d disc%s: ", len(spec.Discs), s))
		for i, disc := range spec.Discs {
			desc := fmt.Sprintf("Disc %d", i)
			if i == 0 {
				desc = "Boot disc"
			}

			output = append(output, fmt.Sprintf("    %s %d GiB, %s grade", desc, disc.Size/1024, disc.StorageGrade))
		}
	} else {
		output = append(output, "No discs specified")
	}
	_, err := wr.Write([]byte(strings.Join(output, "\r\n")))
	return err
}

func FormatAccount(wr io.Writer, a *Account, def *Account, tpl string) error {
	tmpl, err := template.New("accounts").Funcs(templateFuncMap).Funcs(map[string]interface{}{
		"isDefaultAccount": func(a *Account) bool {
			if a.BillingID != 0 && a.BillingID == def.BillingID {
				return true
			}
			return a.Name != "" && a.Name == def.Name
		},
	}).Parse(accountsTemplate + serverTemplate)

	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(wr, tpl, a)
	if err != nil {
		return err
	}

	return nil
}

func FormatOverview(wr io.Writer, accounts []*Account, defaultAccount *Account, username string) error {
	tmpl, err := template.New("accounts").Funcs(templateFuncMap).Funcs(map[string]interface{}{
		"isDefaultAccount": func(a *Account) bool {
			if a.BillingID != 0 && a.BillingID == defaultAccount.BillingID {
				return true
			}
			return a.Name != "" && a.Name == defaultAccount.Name
		},
	}).Parse(accountsTemplate + serverTemplate)
	if err != nil {
		return err
	}
	ownedAccounts := make([]*Account, 0)
	otherAccounts := make([]*Account, 0)
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
	}

	err = tmpl.ExecuteTemplate(wr, "full_overview", data)
	if err != nil {
		return err
	}

	return nil
}
