package brain

import (
	"net"
	"sort"
	"strings"
)

// IPs represent multiple net.IPs
type IPs []*net.IP

// StringSep combines all the IPs into a single string with the given seperator
func (ips IPs) StringSep(sep string) string {
	return strings.Join(ips.Strings(), sep)
}

// Strings returns each IP in this IPs as a string.
func (ips IPs) Strings() (strings []string) {
	strings = make([]string, len(ips))
	for i, ip := range ips {
		strings[i] = ip.String()
	}
	return
}

func (ips IPs) String() string {
	return ips.StringSep(", ")
}

// Sort sorts the given IPs in-place and also returns it, for daisy chaining.
func (ips IPs) Sort() IPs {
	sort.Sort(ips)
	return ips
}

func (ips IPs) Less(i, j int) bool {
	a := *ips[i]
	b := *ips[j]
	if a.Equal(a.To4()) && b.Equal(b.To4()) {
		a4 := a.To4()
		b4 := b.To4()
		for i := 0; i < 4; i++ {
			if a4[i] < b4[i] {
				return true
			} else if a4[i] > b4[i] {
				return false
			}
		}
	} else if a.Equal(a.To16()) && b.Equal(b.To16()) {
		for i := 0; i < 16; i++ {
			if a[i] < b[i] {
				return true
			} else if a[i] > b[i] {
				return false
			}
		}
	}
	if a.Equal(a.To4()) && b.Equal(b.To16()) {
		return true
	}
	return false
}

func (ips IPs) Swap(i, j int) {
	t := ips[i]
	ips[i] = ips[j]
	ips[j] = t
}

func (ips IPs) Len() int {
	return len(ips)
}
