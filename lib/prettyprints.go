package lib

import (
	"fmt"
	"io"
	"math"
	"strings"
	"text/template"
	"unicode"
)

// = name.group (power state) in Zone
//     xxx.yyy.zzz.www n cores, nGiB RAM, nnnGiB storage on n discs
//     discs:
//       + vda - 35GiB, sata grade
//       + vdd - 100GiB, sata grade
//
//     ips: 213.13

const serverTemplate = `{{ define "serversummary" }} ▸ {{.ShortName }} ({{ if .Deleted }}deleted{{ else if .PowerOn }}powered on{{else}}powered off{{end}}) in {{capitalize .ZoneName}}{{ end }}
{{ define "serverspec" }}   {{ .PrimaryIP }} - {{ pluralize "core" "cores" .Cores }}, {{ mibgib .Memory }}, {{.TotalDiscSize "" | gibtib }} on {{ len .Discs | pluralize "disc" "discs"  }}{{ end }}

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

{{ define "servertwoline" }}{{ template "serversummary" . }}
{{ template "serverspec" . }}{{ end }}

{{ define "serverfull" -}}
{{ template "servertwoline" . }}

{{ template "discs" . }}
{{ template "ips" . }}
{{ end }}`

type TemplateChoice string

const (
	OneLine TemplateChoice = "serversummary"
	TwoLine                = "servertwoline"
	All                    = "serverfull"
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

func FormatAccount(wr io.Writer, a *Account) error {
	output := make([]string, 0, 10)

	gs := ""
	if len(a.Groups) != 1 {
		gs = "s"
	}
	ss := ""
	servers := a.CountVirtualMachines()
	if servers != 1 {
		ss = "s"
	}

	groups := make([]string, len(a.Groups))

	for i, g := range a.Groups {
		groups[i] = g.Name
	}
	output = append(output, fmt.Sprintf("%s - Account containing %d server%s across %d group%s", a.Name, servers, ss, len(a.Groups), gs))
	if a.Owner != nil && a.TechnicalContact != nil {
		output = append(output, fmt.Sprintf("Owner: %s %s (%s), Tech Contact: %s %s (%s)", a.Owner.FirstName, a.Owner.LastName, a.Owner.Username, a.TechnicalContact.FirstName, a.TechnicalContact.LastName, a.TechnicalContact.Username))
	}
	output = append(output, "")
	output = append(output, fmt.Sprintf("Groups in this account: %s", strings.Join(groups, ", ")))

	_, err := wr.Write([]byte(strings.Join(output, "\r\n")))
	return err
}
