package testutil

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// DefVM is used to return a default virtual machine for use in testing
var DefVM = lib.VirtualMachineName{Group: "default", Account: "default-account"}

// DefGroup is used to return a default group for use in testing
var DefGroup = lib.GroupName{Group: "default", Account: "default-account"}
