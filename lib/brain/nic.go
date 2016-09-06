package brain

import (
	"net"
)

// NetworkInterface represents a virtual NIC and what IPs it has routed.
type NetworkInterface struct {
	Label string `json:"label"`

	Mac string `json:"mac"`

	// the following can't be set (or at least, so I'm assuming..)

	ID      int `json:"id"`
	VlanNum int `json:"vlan_num"`
	IPs     IPs `json:"ips"`
	// sadly we can't use map[net.IP]*net.IP because net.IP is a slice and slices don't have equality
	// and we can't use map[*net.IP]*net.IP because we could have two identical IPs in different memory locations and they wouldn't be equal. Rubbish.
	ExtraIPs         map[string]*net.IP `json:"extra_ips"`
	VirtualMachineID int                `json:"virtual_machine_id"`
}
