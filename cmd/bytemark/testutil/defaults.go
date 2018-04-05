package testutil

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// DefVM is a default virtual machine for use in testing
var DefVM = lib.VirtualMachineName{Group: "default", Account: "default-account"}

// DefVM is a default group for use in testing
var DefGroup = lib.GroupName{Group: "default", Account: "default-account"}
