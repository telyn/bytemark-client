package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// VirtualMachine represents a VirtualMachine, as passed around from the virtual_machines endpoint
type VirtualMachine struct {
	Autoreboot            bool   `json:"autoreboot_on,omitempty"`
	CdromURL              string `json:"cdrom_url,omitempty"`
	Cores                 int    `json:"cores,omitempty"`
	Memory                int    `json:"memory,omitempty"`
	Name                  string `json:"name,omitempty"`
	PowerOn               bool   `json:"power_on,omitempty"`
	HardwareProfile       string `json:"hardware_profile,omitempty"`
	HardwareProfileLocked bool   `json:"hardware_profile_locked,omitempty"`
	GroupID               int    `json:"group_id,omitempty"`

	// zone name can be set during creation but not changed
	ZoneName string `json:"zone_name,omitempty"`

	// the following cannot be set
	Discs             []*Disc             `json:"discs,omitempty"`
	ID                int                 `json:"id,omitempty"`
	ManagementAddress *net.IP             `json:"management_address,omitempty"`
	Deleted           bool                `json:"deleted,omitempty"`
	Hostname          string              `json:"hostname,omitempty"`
	Head              string              `json:"head,omitempty"`
	NetworkInterfaces []*NetworkInterface `json:"network_interfaces,omitempty"`

	// TODO(telyn): new fields (last_imaged_with and there is another but I forgot)
}

// PrettyPrint outputs a nice human-readable overview of the server to the given writer.
func (vm VirtualMachine) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	const template = `{{ define "server_sgl" }} ▸ {{.ShortName }} ({{ if .Deleted }}deleted{{ else if .PowerOn }}powered on{{else}}powered off{{end}}) in {{capitalize .ZoneName}}{{ end }}
{{ define "server_spec" }}   {{ .PrimaryIP }} - {{ pluralize "core" "cores" .Cores }}, {{ mibgib .Memory }}, {{ if .Discs}}{{.TotalDiscSize "" | gibtib }} on {{ len .Discs | pluralize "disc" "discs"  }}{{ else }}no discs{{ end }}{{ end }}

{{ define "server_discs" -}}
{{- if .Discs }}    discs:
{{- range .Discs }}
      • {{ .Label }} - {{ gibtib .Size }}, {{ .StorageGrade }} grade
{{- end }}

{{ end -}}
{{ end }}

{{ define "server_ips" }}    ips:
{{- range .AllIPv4Addresses.Sort }}
      • {{.}}
{{- end}}
{{- range .AllIPv6Addresses.Sort }}
      • {{.}}
{{- end -}}
{{- end }}

{{ define "server_medium" }}{{ template "server_sgl" . }}
{{ template "server_spec" . }}{{ end }}

{{ define "server_full" -}}
{{ template "server_medium" . }}

{{ template "server_discs" . -}}
{{- template "server_ips" . }}
{{ end }}`
	return prettyprint.Run(wr, template, "server"+string(detail), vm)

}

// TotalDiscSize returns the sum of all disc capacities in the VM for the given storage grade.
// Provide the empty string to sum all discs regardless of storage grade.
func (vm VirtualMachine) TotalDiscSize(storageGrade string) (total int) {
	total = 0
	for _, disc := range vm.Discs {
		if storageGrade == "" || storageGrade == disc.StorageGrade {
			total += disc.Size
		}
	}
	return total
}

// ShortName returns the first two parts of the hostname (i.e. name.group)
func (vm VirtualMachine) ShortName() string {
	bits := strings.SplitN(vm.Hostname, ".", 3)
	return strings.Join(bits[0:2], ".")
}

// FullName returns the first three parts of the hostname (i.e. name.group.account)
func (vm VirtualMachine) FullName() string {
	bits := strings.SplitN(vm.Hostname, ".", 4)
	return strings.Join(bits[0:3], ".")
}

// AllIPv4Addresses flattens all the IPs for a VM into a single IPs (a []*net.IP with some convenience methods)
func (vm VirtualMachine) AllIPv4Addresses() (ips IPs) {
	for _, nic := range vm.NetworkInterfaces {
		for _, ip := range nic.IPs {
			if ip != nil && ip.To4() != nil {
				ips = append(ips, ip)
			}
		}
		for ip := range nic.ExtraIPs {
			netip := net.ParseIP(ip)
			if netip != nil && netip.To4() != nil {
				ips = append(ips, &netip)
			}
		}
	}
	return ips
}

// AllIPv6Addresses flattens all the v6 IPs for a VM into a single IPs (a []*net.IP with some convenience methods)
func (vm VirtualMachine) AllIPv6Addresses() (ips IPs) {
	for _, nic := range vm.NetworkInterfaces {
		for _, ip := range nic.IPs {
			if ip != nil && ip.To4() == nil {
				ips = append(ips, ip)
			}
		}
		for ip := range nic.ExtraIPs {
			netip := net.ParseIP(ip)
			if netip != nil && netip.To4() == nil {
				ips = append(ips, &netip)
			}
		}
	}
	return ips
}

// GetDiscLabelOffset gets the highest disc number for this VM, by looking for discs labelled disc-N and using N or the number of discs attached to the VM, whichever is higher
func (vm VirtualMachine) GetDiscLabelOffset() (offset int) {
	re := regexp.MustCompile(`^dis[ck]-(\d+)$`)
	for _, d := range vm.Discs {
		matches := re.FindStringSubmatch(d.Label)
		if len(matches) < 1 {
			continue
		}
		discNum, err := strconv.ParseInt(matches[1], 10, 32)
		if err != nil {
			discNum = 0
		}
		if int(discNum) > offset {
			offset = int(discNum)
		}
	}
	if offset < len(vm.Discs) {
		return len(vm.Discs)
	}
	return
}

// PrimaryIP returns the VM's primary IP - the (usually) IPv4 address that was created first.
func (vm VirtualMachine) PrimaryIP() net.IP {
	for _, nic := range vm.NetworkInterfaces {
		if len(nic.IPs) > 0 {
			return *nic.IPs[0]
		}
	}
	return nil
}
