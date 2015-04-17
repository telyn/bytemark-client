package cmd

import (
	"testing"
	//"github.com/cheekybits/is"
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

func TestSet(t *testing.T) {
	config := &mockConfig{}
	config.When("SetPersistent", "user", "test-user").Times(1)
	config.When("Get", "user").Return("old-test-user")
	cmds := NewCommandSet(config, nil)
	cmds.Set([]string{"user", "test-user"})

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
}
