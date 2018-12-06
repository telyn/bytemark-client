package flags

import (
	"net"
	"strings"
)

// IPFlag is a flag.Value used to provide an array of net.IPs
type IPFlag []net.IP

// Set sets the IPFlag given the space-separated string of IPs
func (ips *IPFlag) Set(value string) error {
	for _, val := range strings.Split(value, " ") {
		ip := net.ParseIP(val)
		*ips = append(*ips, ip)
	}
	return nil
}

func (ips *IPFlag) String() string {
	var val []string
	for _, ip := range *ips {
		val = append(val, ip.String())
	}
	return strings.Join(val, ", ")
}
