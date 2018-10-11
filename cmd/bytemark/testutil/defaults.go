package testutil

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// DefVM is used to return a default virtual machine for use in testing
var DefVM = lib.VirtualMachineName{Group: "default", Account: "default-account"}

// DefGroup is used to return a default group for use in testing
var DefGroup = lib.GroupName{Group: "default", Account: "default-account"}

// GetFixtureVM returns a default VM for use in testing
func GetFixtureVM() brain.VirtualMachine {
	return brain.VirtualMachine{
		Name:     "test-server",
		Hostname: "test-server.test-group",
		GroupID:  1,
	}
}
