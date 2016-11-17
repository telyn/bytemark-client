package brain

import (
	"net"
)

// IP represents an IP for the purpose of setting RDNS
type IP struct {
	RDns string `json:"rdns"`

	// this cannot be set.
	IP *net.IP `json:"ip"`
}
