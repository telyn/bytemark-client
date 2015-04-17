package cmd

import (
	"testing"
	//"github.com/cheekybits/is"
)

func TestDispatchDoDebug(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Debug", []string{"GET", "/test"})
	d := NewDispatcherWithCommands(config, commands)
	d.Do([]string{"debug", "GET", "/test"})
}

func TestDispatchDoHelp(t *testing.T) {

}

func TestDispatchDoSet(t *testing.T) {

}

func TestDispatchDoShow(t *testing.T) {

}

func TestDispatchDoUnset(t *testing.T) {

}
