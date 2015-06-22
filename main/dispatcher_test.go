package main

import (
	"testing"
	//"github.com/cheekybits/is"
)

func doDispatchTest(t *testing.T, config *mockConfig, commands *mockCommands, args ...string) {
	d := NewDispatcherWithCommands(config, commands)

	if args == nil {
		args = []string{}
	}

	d.Do(args)

	if ok, err := commands.Verify(); !ok {
		t.Fatalf("Test with args %v failed: %v", args, err)
	}
}

func TestDispatchDoCreate(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("CreateGroup", []string{"test-group"}).Times(1)
	doDispatchTest(t, config, commands, "create-group", "test-group")

	commands.Reset()

	commands.When("CreateVM", []string{}).Times(1)
	doDispatchTest(t, config, commands, "create-vm")
}

func TestDispatchDoDebug(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Debug", []string{"GET", "/test"}).Times(1)
	doDispatchTest(t, config, commands, "debug", "GET", "/test")

	commands.Reset()
}

func TestDispatchDoDelete(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Help", []string{"delete-vm"}).Times(1)
	doDispatchTest(t, config, commands, "delete-vm")
	commands.Reset()

	commands.When("DeleteVM", []string{"test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "delete-vm", "test.virtual.machine")
	commands.Reset()
}

func TestDispatchDoUndelete(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Help", []string{"undelete-vm"}).Times(1)
	doDispatchTest(t, config, commands, "undelete-vm")
	commands.Reset()

	commands.When("UndeleteVM", []string{"test.virtual.machine"}).Times(1)
	doDispatchTest(t, config, commands, "undelete-vm", "test.virtual.machine")
	commands.Reset()
}

func TestDispatchDoHelp(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)
	config.When("Get", "silent").Return("true")

	commands.When("Help", []string{}).Times(1)

	doDispatchTest(t, config, commands)

	commands.Reset()
	commands.When("Help", []string{}).Times(1)
	doDispatchTest(t, config, commands, "help")

	commands.When("Help", []string{"show"}).Times(1)
	doDispatchTest(t, config, commands, "help", "show")

	commands.When("Help", []string{"debug"}).Times(1)
	doDispatchTest(t, config, commands, "help", "debug")
}

func TestDispatchDoConfig(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
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

func TestDispatchDoShow(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)
	config.When("Get", "silent").Return(true)

	commands.When("ShowVM", []string{"some-vm"}).Times(1)
	doDispatchTest(t, config, commands, "show-vm", "some-vm")

	commands.Reset()
	commands.When("ShowGroup", []string{"some-group"}).Times(1)
	doDispatchTest(t, config, commands, "show-group", "some-group")

	commands.Reset()
	commands.When("ShowAccount", []string{"some-account"}).Times(1)
	doDispatchTest(t, config, commands, "show-account", "some-account")

}

//func TestDispatchDoUnset(t *testing.T) {
//	commands := &mockCommands{}
//	config := &mockConfig{}
//	config.When("Get", "endpoint").Return("endpoint.example.com")
//	config.When("GetDebugLevel").Return(0)
//
//	doDispatchTest(t, config, commands)
//}
