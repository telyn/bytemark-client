package brain

import (
	"fmt"
	"strings"
)

// IPRange is a representation of an IP range
type IPRange struct {
	ID        int      `json:"id"`
	Spec      string   `json:"spec"`
	VLANNum   int      `json:"vlan_num"`
	Zones     []string `json:"zones"`
	Available float64  `json:"available"` // Needs to be a float64, since the number could go past int64 size
}

// String serialises an IP range to easily be output
func (ipr IPRange) String() string {
	zones := ""
	if len(ipr.Zones) > 0 {
		zones = fmt.Sprintf(", in zones %s", strings.Join(ipr.Zones, ","))
	}
	// Since `Available` is a float64 but won't have decimal points, we just format accordingly.
	return fmt.Sprintf("IP range %s%s, %.0f IPs available.", ipr.Spec, zones, ipr.Available)
}
