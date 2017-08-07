package brain

import (
	"fmt"
	"io"
	"net"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// NetworkInterface represents a virtual NIC and what IPs it has routed.
type NetworkInterface struct {
	Label string `json:"label"`

	Mac string `json:"mac"`

	// the following can't be set (or at least, so I'm assuming..)

	ID      int `json:"id"`
	VlanNum int `json:"vlan_num"`
	IPs     IPs `json:"ips"`
	// sadly we can't use map[net.IP]net.IP because net.IP is a slice and slices don't have equality
	// and we can't use map[*net.IP]net.IP because we could have two identical IPs in different memory locations and they wouldn't be equal. Rubbish.
	ExtraIPs         map[string]net.IP `json:"extra_ips"`
	VirtualMachineID int               `json:"virtual_machine_id"`
}

// ExtraIPStrings represents the extra IPs as a set of strings of the format "extra_ip -> routed_to"
func (nic NetworkInterface) ExtraIPStrings() (ips []string) {
	ips = make([]string, 0, len(nic.ExtraIPs))

	for extraIP, routedTo := range nic.ExtraIPs {
		ips = append(ips, fmt.Sprintf("%s -> %s", extraIP, routedTo))

	}
	return ips
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (nic NetworkInterface) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "Label, Mac, VlanNum"
	}
	return "ID, Label, Mac, VlanNum, IPs, ExtraIPStrings"
}

// PrettyPrint outputs the 
func (nic NetworkInterface) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	nicTpl := `
{{ define "nic_sgl" }}{{ .String }}{{ end }}
{{ define "nic_medium" }}{{ template "nic_sgl" . }}{{ end }}
{{ define "nic_full" -}}
{{- template "nic_medium" . }}
IPs directly attached: {{ join .IPs ", " }}

`

	return prettyprint.Run(wr, nicTpl, "nic"+string(detail), nic)
}

// String formats the network interface as a single descriptive line of text.
func (nic NetworkInterface) String() string {
	return fmt.Sprintf("%s - %s - %d IPs", nic.Label, nic.Mac, len(nic.IPs)+len(nic.ExtraIPs))
}
