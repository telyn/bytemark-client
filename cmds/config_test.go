package cmds

import (
	"bigv.io/client/cmds/util"
	"bigv.io/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestCommandConfig(t *testing.T) {
	config := &mocks.Config{}

	config.When("GetV", "user").Return(util.ConfigVar{"user", "old-test-user", "config"})
	config.When("Get", "user").Return("old-test-user")
	config.When("Silent").Return(true)

	config.When("SetPersistent", "user", "test-user", "CMD set").Times(1)

	cmds := NewCommandSet(config, nil)
	cmds.Config([]string{"set", "user", "test-user"})

	if ok, err := config.Verify(); !ok {
		t.Fatal(err)
	}
}
