package util

import (
	"bytemark.co.uk/client/lib"
	"fmt"
	"net"
	"strings"
)

type IPFlag []net.IP

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

type DiscSpecFlag []lib.Disc

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

func (discFlag *DiscSpecFlag) String() string {
	var discs []string
	for _, d := range *discFlag {
		if d.Label == "" {
			discs = append(discs, fmt.Sprintf("%s:%dGiB", d.StorageGrade, d.Size/1024))
		} else {
			discs = append(discs, fmt.Sprintf("%s:%s:%dGiB", d.Label, d.StorageGrade, d.Size/1024))
		}
	}
	return strings.Join(discs, ",")
}

type SizeSpecFlag int

func (ssf *SizeSpecFlag) Set(spec string) error {
	sz, err := ParseSize(spec)
	if err != nil {
		return err
	}
	*ssf = SizeSpecFlag(sz)
	return nil
}

func (ssf *SizeSpecFlag) String() string {
	return fmt.Sprintf("%d", ssf)
}
