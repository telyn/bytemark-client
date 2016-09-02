package util

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"net"
	"strings"
)

// IPFlag is a flag.Value used to provide an array of net.IPs
type IPFlag []net.IP

// Set sets the IPFlag given the space-seperated string of IPs
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

// DiscSpecFlag is a flag which reads its argument as a disc spec. It can be specified multiple times to add multiple discs.
type DiscSpecFlag []lib.Disc

// Set adds all the defined discs to this flag's value
func (discsFlag *DiscSpecFlag) Set(value string) error {
	for _, val := range strings.Split(value, " ") {
		disc, err := ParseDiscSpec(val)
		if err != nil {
			return err
		}

		*discsFlag = append(*discsFlag, *disc)

	}
	return nil
}

func (discsFlag *DiscSpecFlag) String() string {
	var discs []string
	for _, d := range *discsFlag {
		if d.Label == "" {
			discs = append(discs, fmt.Sprintf("%s:%dGiB", d.StorageGrade, d.Size/1024))
		} else {
			discs = append(discs, fmt.Sprintf("%s:%s:%dGiB", d.Label, d.StorageGrade, d.Size/1024))
		}
	}
	return strings.Join(discs, ",")
}

// SizeSpecFlag represents a capacity as an integer number of megabytes.
type SizeSpecFlag int

// Set sets the value to the size specified. Users can add "M" or "G" as a suffix to specify that they are talking about megabytes/gigabytes. Gigabytes are assumed by default.
func (ssf *SizeSpecFlag) Set(spec string) error {
	sz, err := ParseSize(spec)
	if err != nil {
		return err
	}
	*ssf = SizeSpecFlag(sz)
	return nil
}

func (ssf *SizeSpecFlag) String() string {
	// default value is 0, but is checked for in code that uses SizeSpecFlag and changed to 1
	if *ssf == 0 {
		return ""
	} else {
		return fmt.Sprintf("%d", *ssf)
	}
}
