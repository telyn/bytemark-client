package mocks

import (
	"bigv.io/client/cmd"
	"flag"
	mock "github.com/maraino/go-mock"
)

// mock Config

type mockConfig struct {
	mock.Mock
}

func (c *mockConfig) EndpointName() string {
	ret := c.Called()
	return ret.String(0)
}
func (c *mockConfig) Force() bool {
	ret := c.Called()
	return ret.Bool(0)
}

func (c *mockConfig) Get(name string) (string, error) {
	ret := c.Called(name)
	return ret.String(0), ret.Error(1)
}
func (c *mockConfig) GetIgnoreErr(name string) string {
	ret := c.Called(name)
	return ret.String(0)
}

func (c *mockConfig) GetBool(name string) (bool, error) {
	ret := c.Called(name)
	return ret.Bool(0), ret.Error(1)
}

func (c *mockConfig) GetV(name string) (cmd.ConfigVar, error) {
	ret := c.Called(name)
	return ret.Get(0).(cmd.ConfigVar), ret.Error(1)
}

func (c *mockConfig) GetAll() ([]cmd.ConfigVar, error) {
	ret := c.Called()
	return ret.Get(0).([]cmd.ConfigVar), ret.Error(1)
}

func (c *mockConfig) PanelURL() string {
	ret := c.Called()
	return ret.String(0)
}

func (c *mockConfig) Set(name, value, source string) {
	c.Called(name, value, source)
	return
}

func (c *mockConfig) SetPersistent(name, value, source string) error {
	ret := c.Called(name, value, source)
	return ret.Error(0)
}
func (c *mockConfig) Silent() bool {
	ret := c.Called()
	return ret.Bool(0)
}

func (c *mockConfig) Unset(name string) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *mockConfig) GetDebugLevel() int {
	ret := c.Called()
	return ret.Int(0)
}

func (c *mockConfig) ImportFlags(*flag.FlagSet) []string {
	ret := c.Called()
	if arr, ok := ret.Get(0).([]string); ok {
		return arr
	}
	return nil
}
