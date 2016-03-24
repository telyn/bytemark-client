package util

import (
	"bytemark.co.uk/client/lib"
	"flag"
	"fmt"
	"net"
	"strings"
)

type nilWriter struct{}

func (n *nilWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// MakeCommonFlagSet creates a FlagSet which provides the flags shared between the main command and the sub-commands.
func MakeCommonFlagSet() *flag.FlagSet {
	flags := flag.NewFlagSet("bytemark", flag.ContinueOnError)
	flags.SetOutput(&nilWriter{})

	// because I'm creating my own help functions I'm not going to supply usage info. Neener neener.
	flags.Bool("help", false, "")
	flags.Bool("force", false, "")
	flags.Bool("silent", false, "")
	flags.Bool("yubikey", false, "")
	flags.Int("debug-level", 0, "")
	flags.String("user", "", "")
	flags.String("account", "", "")
	flags.String("endpoint", "", "")
	flags.String("billing-endpoint", "", "")
	flags.String("auth-endpoint", "", "")
	flags.String("config-dir", "", "")
	flags.String("yubikey-otp", "", "")

	return flags
}

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
