package brain

import (
	"net"
)

// IPCreateRequest is used by the create_ip endpoint on a nic
type IPCreateRequest struct {
	Addresses  int    `json:"addresses"`
	Family     string `json:"family"`
	Reason     string `json:"reason"`
	Contiguous bool   `json:"contiguous"`
	// don't actually specify the IPs - this is for filling in from the response!
	IPs []*net.IP `json:"ips"`
}
