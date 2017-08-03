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

func (ip IP) DefaultFields(f output.Format) string {
	return "IP, RDns"
}

func (ip IP) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	fmt.Fprintf(wr, "%s: %s", ip.IP, ip.RDns)
	return nil
}
