package brain

import (
	"fmt"
	"io"
	"net"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// IP represents an IP for the purpose of setting RDNS
type IP struct {
	RDns string `json:"rdns"`

	// this cannot be set.
	IP net.IP `json:"ip"`
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/bytemark-client for this type.
func (ip IP) DefaultFields(f output.Format) string {
	return "IP, RDns"
}

// PrettyPrint outputs a vaguely human-readable version of the IP and reverse DNS to wr. Detail is ignored.
func (ip IP) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	fmt.Fprintf(wr, "%s: %s", ip.IP, ip.RDns)
	return nil
}
