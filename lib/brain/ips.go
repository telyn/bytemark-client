package brain

import (
	"net"
	"sort"
	"strings"
)

// IPs represent multiple net.IPs
type IPs []net.IP

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
// Sorting is done by the golang sort package, which uses the Less, Len and Swap functions defined below
func (ips IPs) Sort() IPs {
	sort.Sort(ips)
	return ips
}

// Less looks at the ips at index i and j, and returns true if i should come before j.
func (ips IPs) Less(i, j int) bool {
	a := ips[i]
	b := ips[j]
	// loop over each byte of the address and compare.
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
	// v4 < v6 always
	if a.Equal(a.To4()) && b.Equal(b.To16()) {
		return true
	}
	return false
}

// Swap moves the ip at i to j, and vice versa.
func (ips IPs) Swap(i, j int) {
	t := ips[i]
	ips[i] = ips[j]
	ips[j] = t
}

// Len returns how many ips there are in this IPs
func (ips IPs) Len() int {
	return len(ips)
}
