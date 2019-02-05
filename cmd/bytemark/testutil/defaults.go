package testutil

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// DefVM is used to return a default virtual machine for use in testing
var DefVM = pathers.VirtualMachineName{GroupName: pathers.GroupName{Group: "default", Account: "default-account"}}

// DefGroup is used to return a default group for use in testing
var DefGroup = pathers.GroupName{Group: "default", Account: "default-account"}

// GetFixtureVM returns a default VM for use in testing
func GetFixtureVM() brain.VirtualMachine {
	return brain.VirtualMachine{
		Name:     "test-server",
		Hostname: "test-server.test-group",
		GroupID:  1,
	}
}
