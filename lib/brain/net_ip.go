package brain

import (
	"fmt"
	"net"
)

// NetIP is an alias to *net.IP which implements lib.Pather
type NetIP net.IP

// Path returns a URL path to access a virtual_machine for the given ip
func (ip NetIP) Path() (string, error) {
	return fmt.Sprintf("/virtual_machines/%s", (net.IP)(ip).String()), nil
}
