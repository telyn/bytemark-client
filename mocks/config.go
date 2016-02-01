package mocks

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/lib"
	"flag"
	mock "github.com/maraino/go-mock"
)

// mock Config

type Config struct {
	mock.Mock
}

func (c *Config) EndpointName() string {
	ret := c.Called()
	return ret.String(0)
}
func (c *Config) Force() bool {
	ret := c.Called()
	return ret.Bool(0)
}

func (c *Config) Get(name string) (string, error) {
	ret := c.Called(name)
	return ret.String(0), ret.Error(1)
}
func (c *Config) GetIgnoreErr(name string) string {
	ret := c.Called(name)
	return ret.String(0)
}

func (c *Config) GetBool(name string) (bool, error) {
	ret := c.Called(name)
	return ret.Bool(0), ret.Error(1)
}

func (c *Config) GetV(name string) (util.ConfigVar, error) {
	ret := c.Called(name)
	return ret.Get(0).(util.ConfigVar), ret.Error(1)
}

func (c *Config) GetVirtualMachine() lib.VirtualMachineName {
	ret := c.Called()
	return ret.Get(0).(lib.VirtualMachineName)
}

func (c *Config) GetGroup() lib.GroupName {
	ret := c.Called()
	return ret.Get(0).(lib.GroupName)
}

func (c *Config) GetAll() ([]util.ConfigVar, error) {
	ret := c.Called()
	return ret.Get(0).([]util.ConfigVar), ret.Error(1)
}

func (c *Config) PanelURL() string {
	ret := c.Called()
	return ret.String(0)
}

func (c *Config) Set(name, value, source string) {
	c.Called(name, value, source)
	return
}

func (c *Config) SetPersistent(name, value, source string) error {
	ret := c.Called(name, value, source)
	return ret.Error(0)
}
func (c *Config) Silent() bool {
	ret := c.Called()
	return ret.Bool(0)
}

func (c *Config) Unset(name string) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *Config) GetDebugLevel() int {
	ret := c.Called()
	return ret.Int(0)
}

func (c *Config) ImportFlags(*flag.FlagSet) []string {
	ret := c.Called()
	if arr, ok := ret.Get(0).([]string); ok {
		return arr
	}
	return nil
}
