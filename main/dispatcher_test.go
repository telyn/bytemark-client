package main

import (
	"bigv.io/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func doDispatchTest(t *testing.T, config *mocks.Config, commands *mocks.Commands, args ...string) {
	d, err := NewDispatcherWithCommandManager(config, commands)
	if err != nil {
		t.Fatalf("NewDispatcherWithCommands died: %v", err)
	}

	if args == nil {
		args = []string{}
	}

	d.Do(args)

	if ok, err := commands.Verify(); !ok {
		t.Fatalf("Test with args %v failed: %v", args, err)
	}
}

func TestDispatchDoCreate(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("CreateGroup", []string{"test-group"}).Times(1)
	doDispatchTest(t, config, commands, "create", "group", "test-group")

	commands.Reset()

	commands.When("CreateVM", []string{}).Times(1)
	doDispatchTest(t, config, commands, "create", "vm")

	commands.When("CreateDiscs", []string{}).Times(1)
	doDispatchTest(t, config, commands, "create", "disc")

	commands.When("CreateDiscs", []string{}).Times(1)
	doDispatchTest(t, config, commands, "create", "discs")

	commands.When("CreateDiscs", []string{}).Times(1)
	doDispatchTest(t, config, commands, "create", "disk")
	commands.When("CreateDiscs", []string{}).Times(1)
	doDispatchTest(t, config, commands, "create", "disks")
}

func TestDispatchDoDebug(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Debug", []string{"GET", "/test"}).Times(1)
	doDispatchTest(t, config, commands, "debug", "GET", "/test")

	commands.Reset()
}

func TestDispatchDoDelete(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("DeleteVM", []string{}).Times(1)
	doDispatchTest(t, config, commands, "delete", "vm")
	commands.Reset()

	commands.When("DeleteVM", []string{"test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "delete", "vm", "test.virtual.machine")
	commands.Reset()

	commands.When("DeleteGroup", []string{"test-group.account"}).Times(1)
	doDispatchTest(t, config, commands, "delete", "group", "test-group.account")
	commands.Reset()
}

func TestDispatchDoUndelete(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("UndeleteVM", []string{}).Times(1)
	doDispatchTest(t, config, commands, "undelete", "vm")
	commands.Reset()

	commands.When("UndeleteVM", []string{"test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "undelete", "vm", "test.virtual.machine")
	commands.Reset()
}

func TestDispatchDoLock(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("HelpForLocks").Times(1)
	doDispatchTest(t, config, commands, "lock")
	commands.Reset()

	commands.When("HelpForLocks").Times(1)
	doDispatchTest(t, config, commands, "lock", "non-existent")
	commands.Reset()

	commands.When("LockHWProfile", []string{}).Times(1)
	doDispatchTest(t, config, commands, "lock", "hwprofile")

	commands.When("HelpForLocks").Times(0)
	commands.When("LockHWProfile", []string{"test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "lock", "hwprofile", "test.virtual.machine")
	commands.Reset()
}

func TestDispatchDoUnlock(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("HelpForLocks").Times(1)
	doDispatchTest(t, config, commands, "unlock")
	commands.Reset()

	commands.When("HelpForLocks").Times(1)
	doDispatchTest(t, config, commands, "unlock", "non-existent")
	commands.Reset()

	commands.When("UnlockHWProfile", []string{}).Times(1)
	doDispatchTest(t, config, commands, "unlock", "hwprofile")
	commands.Reset()

	commands.When("HelpForLocks").Times(0)
	commands.When("UnlockHWProfile", []string{"test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "unlock", "hwprofile", "test.virtual.machine")
	commands.Reset()
}

func TestDispatchDoSet(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("HelpForSet").Times(1)
	doDispatchTest(t, config, commands, "set")
	commands.Reset()

	// currently only the dispatcher is tested, HelpForSet IS called but in the
	// command which isn't loaded, so these tests are disable for now
	//commands.When("HelpForSet").Times(1)
	//doDispatchTest(t, config, commands, "set", "hwprofile")
	//commands.Reset()

	//commands.When("HelpForSet").Times(1)
	//doDispatchTest(t, config, commands, "set", "hwprofile", "test.virtual.machine")
	//commands.Reset()

	commands.When("HelpForSet").Times(0)
	commands.When("SetHWProfile", []string{"--locked", "test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "set", "hwprofile", "--locked", "test.virtual.machine")
	commands.Reset()

	commands.When("HelpForSet").Times(0)
	commands.When("SetHWProfile", []string{"--unlocked", "test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "set", "hwprofile", "--unlocked", "test.virtual.machine")
	commands.Reset()

	commands.When("HelpForSet").Times(0)
	commands.When("SetHWProfile", []string{"test.virtual.machine", "virtio123"}).Times(1)
	doDispatchTest(t, config, commands, "set", "hwprofile", "test.virtual.machine", "virtio123")
	commands.Reset()

	commands.When("HelpForSet").Times(0)
	commands.When("SetCores", []string{"test.virtual.machine", "10"}).Times(1)
	doDispatchTest(t, config, commands, "set", "cores", "test.virtual.machine", "10")
	commands.Reset()

	commands.When("HelpForSet").Times(0)
	commands.When("SetCores", []string{"test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "set", "cores", "test.virtual.machine")
	commands.Reset()

	commands.When("HelpForSet").Times(0)
	commands.When("SetMemory", []string{"test.virtual.machine", "16G"}).Times(1)
	doDispatchTest(t, config, commands, "set", "memory", "test.virtual.machine", "16G")
	commands.Reset()

	commands.When("HelpForSet").Times(0)
	commands.When("SetMemory", []string{"test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "set", "memory", "test.virtual.machine")
	commands.Reset()
}

func TestDispatchDoHelp(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)
	config.When("Get", "silent").Return("true")

	commands.When("Help", []string{}).Times(1)
	doDispatchTest(t, config, commands)

	commands.When("Help", []string{}).Times(1)
	doDispatchTest(t, config, commands, "help")

	commands.When("Help", []string{"show"}).Times(1)
	doDispatchTest(t, config, commands, "help", "show")

	commands.When("Help", []string{"debug"}).Times(1)
	doDispatchTest(t, config, commands, "help", "debug")
}

func TestDispatchDoConfig(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Config", []string{}).Times(1)
	doDispatchTest(t, config, commands, "config")

	commands.Reset()
	commands.When("Config", []string{"set"}).Times(1)
	doDispatchTest(t, config, commands, "config", "set")

	commands.Reset()
	commands.When("Config", []string{"set", "variablename"}).Times(1)
	doDispatchTest(t, config, commands, "config", "set", "variablename")

	commands.Reset()
	commands.When("Config", []string{"set", "variablename", "value"}).Times(1)
	doDispatchTest(t, config, commands, "config", "set", "variablename", "value")
}

func TestDispatchDoPower(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Shutdown", []string{}).Times(1)
	doDispatchTest(t, config, commands, "shutdown")

	commands.Reset()
	commands.When("Start", []string{}).Times(1)
	doDispatchTest(t, config, commands, "start")

	commands.Reset()
	commands.When("Stop", []string{}).Times(1)
	doDispatchTest(t, config, commands, "stop")

	commands.Reset()
	commands.When("Restart", []string{}).Times(1)
	doDispatchTest(t, config, commands, "restart")

	commands.Reset()
	commands.When("ResetVM", []string{}).Times(1)
	doDispatchTest(t, config, commands, "reset")

	commands.Reset()
	commands.When("Shutdown", []string{"test-vm"}).Times(1)
	doDispatchTest(t, config, commands, "shutdown", "test-vm")

	commands.Reset()
	commands.When("Start", []string{"test-vm"}).Times(1)
	doDispatchTest(t, config, commands, "start", "test-vm")

	commands.Reset()
	commands.When("Stop", []string{"test-vm"}).Times(1)
	doDispatchTest(t, config, commands, "stop", "test-vm")

	commands.Reset()
	commands.When("Restart", []string{"test-vm"}).Times(1)
	doDispatchTest(t, config, commands, "restart", "test-vm")

	commands.Reset()
	commands.When("ResetVM", []string{"test-vm"}).Times(1)
	doDispatchTest(t, config, commands, "reset", "test-vm")

}

func TestDispatchDoShow(t *testing.T) {
	commands := &mocks.Commands{}
	config := &mocks.Config{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)
	config.When("Get", "silent").Return(true)

	commands.When("ShowVM", []string{"some-vm"}).Times(1)
	doDispatchTest(t, config, commands, "show", "vm", "some-vm")

	commands.Reset()
	commands.When("ShowGroup", []string{"some-group"}).Times(1)
	doDispatchTest(t, config, commands, "show", "group", "some-group")

	commands.Reset()
	commands.When("ShowAccount", []string{"some-account"}).Times(1)
	doDispatchTest(t, config, commands, "show", "account", "some-account")

}

//func TestDispatchDoUnset(t *testing.T) {
//	commands := &mocks.Commands{}
//	config := &mocks.Config{}
//	config.When("Get", "endpoint").Return("endpoint.example.com")
//	config.When("GetDebugLevel").Return(0)
//
//	doDispatchTest(t, config, commands)
//}
