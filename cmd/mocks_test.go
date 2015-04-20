package cmd

import (
	mock "github.com/maraino/go-mock"
)

type mockConfig struct {
	mock.Mock
}

func (c *mockConfig) Get(name string) string {
	ret := c.Called(name)
	return ret.String(0)
}

func (c *mockConfig) GetV(name string) ConfigVar {
	ret := c.Called(name)
	return ret.Get(0).(ConfigVar)
}

func (c *mockConfig) GetAll() []ConfigVar {
	ret := c.Called()
	return ret.Get(0).([]ConfigVar)
}

func (c *mockConfig) Set(name, value, source string) {
	c.Called(name, value, source)
	return
}

func (c *mockConfig) SetPersistent(name, value, source string) {
	c.Called(name, value, source)
	return
}

func (c *mockConfig) Unset(name string) {
	c.Called(name)
	return
}

func (c *mockConfig) GetDebugLevel() int {
	ret := c.Called()
	return ret.Int(0)
}

type mockCommands struct {
	mock.Mock
}

func (cmds *mockCommands) Debug(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) Help(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) Config(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) ShowAccount(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) ShowVM(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) EnsureAuth() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForConfig() {
	cmds.Called()
}
func (cmds *mockCommands) HelpForDebug() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForHelp() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForShow() {
	cmds.Called()
}
