package flags

import (
	"net"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// AccountName returns the named AccountNameFlag. Why does it return the
// flag and not simply the account name? I don't know
func AccountName(c *app.Context, flagname string) AccountNameFlag {
	accountName, ok := c.Context.Generic(flagname).(*AccountNameFlag)
	if ok {
		return AccountNameFlag(*accountName)
	}
	return AccountNameFlag{}
}

// Discs returns the discs passed along as the named flag.
// I can't imagine why I'd ever name a disc flag anything other than --disc, but the flexibility is there just in case.
func Discs(c *app.Context, flagname string) []brain.Disc {
	disc, ok := c.Context.Generic(flagname).(*DiscSpecFlag)
	if ok {
		return []brain.Disc(*disc)
	}
	return []brain.Disc{}
}

// FileName returns the name of the file given by the named flag
func FileName(c *app.Context, flagname string) string {
	file, ok := c.Context.Generic(flagname).(*FileFlag)
	if ok {
		return file.FileName
	}
	return ""
}

// FileContents returns the contents of the file given by the named flag.
func FileContents(c *app.Context, flagname string) string {
	file, ok := c.Context.Generic(flagname).(*FileFlag)
	if ok {
		return file.Value
	}
	return ""
}

// GroupName returns the named flag as a pathers.GroupName
func GroupName(c *app.Context, flagname string) (gp pathers.GroupName) {
	gpNameFlag, ok := c.Context.Generic(flagname).(*GroupNameFlag)
	if !ok {
		return pathers.GroupName{}
	}
	if gpNameFlag == nil {
		return pathers.GroupName{}
	}
	return gpNameFlag.GroupName
}

// IPs returns the ips passed along as the named flag.
func IPs(c *app.Context, flagname string) []net.IP {
	ips, ok := c.Context.Generic(flagname).(*IPFlag)
	if ok {
		return []net.IP(*ips)
	}
	return []net.IP{}
}

// Privilege returns the named flag as a PrivilegeFlag
func Privilege(c *app.Context, flagname string) PrivilegeFlag {
	priv, ok := c.Context.Generic(flagname).(*PrivilegeFlag)
	if ok {
		return *priv
	}
	return PrivilegeFlag{}
}

// Resize returns the named ResizeFlag
func Resize(c *app.Context, flagname string) ResizeFlag {
	size, ok := c.Context.Generic(flagname).(*ResizeFlag)
	if ok {
		return *size
	}
	return ResizeFlag{}
}

// Size returns the value of the named SizeSpecFlag as an int in megabytes
func Size(c *app.Context, flagname string) int {
	size, ok := c.Context.Generic(flagname).(*SizeSpecFlag)
	if ok {
		return int(*size)
	}
	return 0
}

// VirtualMachineName returns the named flag as a lib.VirtualMachineName
func VirtualMachineName(c *app.Context, flagname string) (vm lib.VirtualMachineName) {
	vmNameFlag, ok := c.Context.Generic(flagname).(*VirtualMachineNameFlag)
	if !ok {
		return c.Config().GetVirtualMachine()
	}
	if vmNameFlag == nil {
		return lib.VirtualMachineName{}
	}

	return vmNameFlag.VirtualMachineName
}
