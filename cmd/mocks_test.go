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

func (c *mockConfig) Set(name, value string) {
	c.Called(name, value)
	return
}

func (c *mockConfig) SetPersistent(name, value string) {
	c.Called(name, value)
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

func (cmds *mockCommands) Set(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) ShowAccount(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) ShowVM(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) Unset(args []string) {
	cmds.Called(args)
}

func (cmds *mockCommands) EnsureAuth() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForDebug() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForHelp() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForSet() {
	cmds.Called()
}

func (cmds *mockCommands) HelpForShow() {
	cmds.Called()
}
