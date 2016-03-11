package main

import (
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/mocks"
	"testing"
)

func TestResetCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-server.test-group.test-account"})
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []lib.VirtualMachineName{{}}).Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("ResetVirtualMachine", vmn).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.ResetServer([]string{"test-server.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestRestartCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-server.test-group.test-account"})
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []lib.VirtualMachineName{{}}).Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("RestartVirtualMachine", vmn).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.Restart([]string{"test-server.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestShutdownCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-server.test-group.test-account"})
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []lib.VirtualMachineName{{}}).Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	c.When("ShutdownVirtualMachine", vmn, true).Times(1)
	cmds.Shutdown([]string{"test-server.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStartCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-server.test-group.test-account"})
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []lib.VirtualMachineName{{}}).Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	cmds := NewCommandSet(config, c)

	c.When("StartVirtualMachine", vmn).Times(1)
	cmds.Start([]string{"test-server.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStopCommand(t *testing.T) {
	c := &mocks.Client{}
	config := &mocks.Config{}

	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("Get", "token").Return("test-token")
	config.When("Silent").Return(true)
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("ImportFlags").Return([]string{"test-server.test-group.test-account"})
	config.When("GetVirtualMachine").Return(lib.VirtualMachineName{})

	c.When("ParseVirtualMachineName", "test-server.test-group.test-account", []lib.VirtualMachineName{{}}).Return(vmn)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	c.When("StopVirtualMachine", vmn).Times(1)

	cmds := NewCommandSet(config, c)

	cmds.Stop([]string{"test-server.test-group.test-account"})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
