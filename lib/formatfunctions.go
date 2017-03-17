package lib

// This file is a compatibility shim for keeping Format* functions (approximately) working before their removal in bytemark-client 3.0
// do not rely on its contents continued existence.

// TODO(telyn): delete this file in 3.0

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"io"
)

// TemplateChoice is which template to use for a server.
// This type is deprecated and will be removed in bytemark-client 3.0
type TemplateChoice string

func translateTemplateChoice(name TemplateChoice) prettyprint.DetailLevel {
	switch name {
	case "server_spec", "server_twoline":
		return prettyprint.Medium
	case "server_name", "account_name", "server_oneline", "account_bullet", "server_summary":
		return prettyprint.SingleLine
	}
	return prettyprint.Full
}

// FormatVirtualMachine outputs the given VM to the given Writer, using the template specified.
// This function is deprecated and will be removed in 3.0. Please use brain.VirtualMachine.PrettyPrint(io.Writer, prettyprint.Detail) instead.
func FormatVirtualMachine(wr io.Writer, vm *brain.VirtualMachine, tpl TemplateChoice) error {
	detail := translateTemplateChoice(TemplateChoice(tpl))

	return vm.PrettyPrint(wr, detail)
}

// FormatImageInstall outputs the given ImageInstall to the given Writer. tpl is ignored.
// This function is deprecated and will be removed in 3.0. Please use brain.ImageInstall.PrettyPrint(io.Writer, prettyprint.Detail) instead.
func FormatImageInstall(wr io.Writer, ii *brain.ImageInstall, tpl TemplateChoice) error {
	return ii.PrettyPrint(wr, prettyprint.Full)
}

// FormatVirtualMachineSpec outputs the given spec to the given writer. tpl and group are ignored.
// This function is deprecated and will be removed in 3.0. Please use brain.VirtualMachineSpec.PrettyPrint(io.Writer, prettyprint.Detail) instead.
func FormatVirtualMachineSpec(wr io.Writer, group *GroupName, spec *brain.VirtualMachineSpec, tpl TemplateChoice) error {
	return spec.PrettyPrint(wr, prettyprint.Full)
}

// FormatAccount outputs the given Account to the given Writer.
// This function is deprecated and will be removed in 3.0. Please use Account.PrettyPrint(io.Writer, prettyprint.Detail) instead.
func FormatAccount(wr io.Writer, a *Account, def *Account, tpl string) error {
	detail := translateTemplateChoice(TemplateChoice(tpl))
	if def != nil && (a.Name == def.Name || a.BillingID == def.BillingID) {
		a.IsDefaultAccount = true
	}
	return a.PrettyPrint(wr, detail)
}
