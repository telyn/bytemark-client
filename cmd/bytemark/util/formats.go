package util

import (
	"fmt"
	client "github.com/BytemarkHosting/bytemark-client/lib"
	"math"
	"strings"
)

// VMFormatOptions controls formatting of VMs in FormatVirtualMachine
// Add or or them together to get what you want
type VMFormatOptions uint8

const (
	// _FormatVMWithAddrs causes IP addresses to be included in the output
	_FormatVMWithAddrs VMFormatOptions = 1 << iota
	// _FormatVMWithDiscs causes individual disc sizes & storage grades to be included in the output
	_FormatVMWithDiscs
	// _FormatVMWithCDURL causes the URL of the image being used as the CD to be included in the output, if applicable
	_FormatVMWithCDURL
	// _FormatVMSingleLine causes a minimal set of details to be output on a single line, and overrides the other options
	_FormatVMSingleLine
	// _FormatVMTwoLine causes a relatively minimal set of details to be output on two lines and overrides the other options with the exception of _FormatVMSingleLine
	_FormatVMTwoLine
)

// FORMAT_DEFAULT_WIDTH is the default width to attempt to print to.
const _FormatDefaultWidth = 80

// FlowString adds newlines and spaces such that the string never gets longer than width, and remains consistently indented to the given number of spaces. It is dumb and assumes every character has a width of one.
// Really it should do it by words, but that's harder and I just want *something* right now
func FlowString(width int, indent int, str string) string {
	lines := make([]string, 0)
	curline := make([]rune, 0, width)

	for i := 0; i < indent; i++ {
		curline = append(curline, ' ')
	}
	ch := 0
	for _, char := range str {
		if ch == width {
			lines = append(lines, string(curline))
			curline = make([]rune, 0, width)
			ch = indent
			for i := 0; i < indent; i++ {
				curline = append(curline, ' ')
			}
		}
		ch++
		curline = append(curline, char)
	}
	lines = append(lines, string(curline))
	return strings.Join(lines, "\r\n")
}

// FlowStringf formats the given string like Sprintf then runs FlowString on the result.
func FlowStringf(width int, indent int, str string, params ...interface{}) string {
	str = fmt.Sprintf(str, params...)
	return FlowString(width, indent, str)
}

// FormatVirtualMachines loops through a bunch of VMs, formatting each one as it goes, and returns each formatted VM as a string.
// The options are the same as FormatVirtualMachine
func FormatVirtualMachines(vms []*client.VirtualMachine, options ...int) []string {
	output := make([]string, len(vms), len(vms))
	for i, vm := range vms {
		output[i] = FormatVirtualMachine(vm, options...)
	}
	return output
}

// FormatVirtualMachine pretty-prints a VM. The optional second argument is a bitmask of VMFormatOptions,
// and the optional third is the width you'd like to display..oh.
func FormatVirtualMachine(vm *client.VirtualMachine, options ...int) string {
	width := _FormatDefaultWidth
	format := _FormatVMWithAddrs | _FormatVMWithDiscs

	if len(options) >= 1 {
		format = VMFormatOptions(options[0])
	}

	if len(options) >= 2 {
		width = options[1]
	}

	output := make([]string, 0, 10)

	// outputIndentedf formats the given string, flows it so it's indented two spaces, and appends it to the output array.
	outputIndentedf := func(str string, params ...interface{}) {
		output = append(output, FlowStringf(width, 4, str, params...))
	}

	powerstate := "powered off"
	if vm.PowerOn {
		powerstate = "powered on"
	}
	if vm.Deleted {
		powerstate = "deleted"
	}

	hostnameparts := strings.Split(vm.Hostname, ".")
	shortname := hostnameparts[0] + "." + hostnameparts[1]
	zone := strings.Title(vm.ZoneName)

	// append & format by hand because we dont want this to be indented
	output = append(output, fmt.Sprintf("= %s (%s) in %s", shortname, powerstate, zone))

	// if !singleline
	if format&_FormatVMSingleLine == _FormatVMSingleLine {
		return output[0]
	}

	memAmt := vm.Memory
	memCh := 'M'
	if memAmt >= 1024 {
		memAmt /= 1024
		memCh = 'G'
	}

	diskAmt := float64(vm.TotalDiscSize("")) / 1024
	diskCh := 'G'
	if diskAmt >= 1024 {
		diskCh = 'T'
		diskAmt = diskAmt / 1024
	}
	diskN := len(vm.Discs)
	sForDisks := "s"
	sForCores := "s"
	if vm.Cores == 1 {
		sForCores = ""
	}
	if diskN == 1 {
		sForDisks = ""
	}

	if diskCh == 'T' {
		outputIndentedf("%d core%s, %d%ciB RAM, %.1f%ciB storage on %d disk%s", vm.Cores, sForCores, memAmt, memCh, diskAmt, diskCh, diskN, sForDisks)
	} else {
		outputIndentedf("%d core%s, %d%ciB RAM, %.0f%ciB storage on %d disk%s", vm.Cores, sForCores, memAmt, memCh, diskAmt, diskCh, diskN, sForDisks)
	}

	if format&_FormatVMTwoLine == _FormatVMTwoLine {
		return strings.Join(output, "\r\n")
	}
	output = append(output, "")

	if (format&_FormatVMWithCDURL) != 0 && vm.CdromURL != "" {
		outputIndentedf("CD-ROM: %s", vm.CdromURL)
	}

	output = append(output, "")
	if (format & _FormatVMWithDiscs) != 0 {
		for _, disc := range vm.Discs {
			outputIndentedf("Disc %s: %d GiB, %s grade", disc.Label, disc.Size/1024, disc.StorageGrade)
		}
		output = append(output, "")
	}

	if (format & _FormatVMWithAddrs) != 0 {
		outputIndentedf("IPv4 Addresses: %s", vm.AllIPv4Addresses().StringSep(",\r\n                "))
		outputIndentedf("IPv6 Addresses: %s", vm.AllIPv6Addresses().StringSep(",\r\n                "))
	}

	return strings.Join(output, "\r\n")
}

func FormatVirtualMachineSpec(group *client.GroupName, spec *client.VirtualMachineSpec) string {
	output := make([]string, 0, 10)
	output = append(output, fmt.Sprintf("Name: '%s'", spec.VirtualMachine.Name))
	output = append(output, fmt.Sprintf("Group: '%s'", group.Group))
	if group.Account == "" {
		output = append(output, "Account: not specified - will default to the account with the same name as the user you log in as")
	} else {
		output = append(output, fmt.Sprintf("Account: '%s'", group.Account))
	}
	s := ""
	if spec.VirtualMachine.Cores > 1 {
		s = "s"
	}

	mems := fmt.Sprintf("%d", spec.VirtualMachine.Memory/1024)
	if 0 != math.Mod(float64(spec.VirtualMachine.Memory), 1024) {
		mem := float64(spec.VirtualMachine.Memory) / 1024.0
		mems = fmt.Sprintf("%.2f", mem)
	}
	output = append(output, fmt.Sprintf("Specs: %d core%s and %sGiB memory", spec.VirtualMachine.Cores, s, mems))

	locked := ""
	if spec.VirtualMachine.HardwareProfile != "" {
		if spec.VirtualMachine.HardwareProfileLocked {
			locked = " (locked)"
		}
		output = append(output, fmt.Sprintf("Hardware profile: %s%s", spec.VirtualMachine.HardwareProfile, locked))
	}

	if spec.IPs != nil {
		if spec.IPs.IPv4 != "" {
			output = append(output, fmt.Sprintf("IPv4 address: %s", spec.IPs.IPv4))
		}
		if spec.IPs.IPv6 != "" {
			output = append(output, fmt.Sprintf("IPv6 address: %s", spec.IPs.IPv6))
		}
	}

	if spec.Reimage != nil {
		if spec.Reimage.Distribution == "" {
			if spec.VirtualMachine.CdromURL == "" {
				output = append(output, "No image or CD URL specified")
			} else {
				output = append(output, fmt.Sprintf("CD URL: %s", spec.VirtualMachine.CdromURL))
			}
		} else {
			output = append(output, "Image: "+spec.Reimage.Distribution)
		}
		output = append(output, "Root/Administrator password: "+spec.Reimage.RootPassword)
	} else {

		if spec.VirtualMachine.CdromURL == "" {
			output = append(output, "No image or CD URL specified")
		} else {
			output = append(output, fmt.Sprintf("CD URL: %s", spec.VirtualMachine.CdromURL))
		}
	}

	s = ""
	if len(spec.Discs) > 1 {
		s = "s"
	}
	if len(spec.Discs) > 0 {
		output = append(output, fmt.Sprintf("%d disc%s: ", len(spec.Discs), s))
		for i, disc := range spec.Discs {
			desc := fmt.Sprintf("Disc %d", i)
			if i == 0 {
				desc = "Boot disc"
			}

			output = append(output, fmt.Sprintf("    %s %d GiB, %s grade", desc, disc.Size/1024, disc.StorageGrade))
		}
	} else {
		output = append(output, "No discs specified")
	}
	return strings.Join(output, "\r\n")

}

func FormatImageInstall(ii *client.ImageInstall) string {
	output := make([]string, 0)
	if ii.Distribution != "" {
		output = append(output, "Image: "+ii.Distribution)
	}
	if ii.PublicKeys != "" {
		keynames := make([]string, 0)
		for _, k := range strings.Split(ii.PublicKeys, "\n") {
			kbits := strings.SplitN(k, " ", 3)
			if len(kbits) == 3 {
				keynames = append(keynames, strings.TrimSpace(kbits[2]))
			}

		}
		output = append(output, fmt.Sprintf("%d public keys: %s", len(keynames), strings.Join(keynames, ", ")))
	}
	if ii.RootPassword != "" {
		output = append(output, "Root/Administrator password: "+ii.RootPassword)
	}
	if ii.FirstbootScript != "" {
		output = append(output, "With a firstboot script")
	}
	return strings.Join(output, "\r\n")
}

func FormatAccount(a *client.Account) string {
	output := make([]string, 0, 10)

	gs := ""
	if len(a.Groups) != 1 {
		gs = "s"
	}
	ss := ""
	servers := a.CountVirtualMachines()
	if servers != 1 {
		ss = "s"
	}

	groups := make([]string, len(a.Groups))

	for i, g := range a.Groups {
		groups[i] = g.Name
	}
	output = append(output, fmt.Sprintf("%s - Account containing %d server%s across %d group%s", a.Name, servers, ss, len(a.Groups), gs))
	if a.Owner != nil && a.TechnicalContact != nil {
		output = append(output, fmt.Sprintf("Owner: %s %s (%s), Tech Contact: %s %s (%s)", a.Owner.FirstName, a.Owner.LastName, a.Owner.Username, a.TechnicalContact.FirstName, a.TechnicalContact.LastName, a.TechnicalContact.Username))
	}
	output = append(output, "")
	output = append(output, fmt.Sprintf("Groups in this account: %s", strings.Join(groups, ", ")))

	return strings.Join(output, "\r\n")

}
