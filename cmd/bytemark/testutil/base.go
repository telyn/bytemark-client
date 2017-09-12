package testutil

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func BaseTestSetup(t *testing.T, admin bool, commands []cli.Command) (config *mocks.Config, client *mocks.Client, cliapp *cli.App) {
	config = new(mocks.Config)
	client = new(mocks.Client)
	config.When("GetBool", "admin").Return(admin, nil)
	config.When("GetV", "output-format").Return(util.ConfigVar{"output-format", "human", "CODE"})

	cliapp, err := app.BaseAppSetup(app.GlobalFlags(), config, commands)
	if err != nil {
		t.Fatal(err)
	}
	cliapp.Metadata = map[string]interface{}{
		"client": client,
		"config": config,
	}

	cliapp.Writer = ioutil.Discard
	for _, c := range commands {
		//config.When("Get", "token").Return("no-not-a-token")

		// the issue is that Command.FullName() is dependent on Command.commandNamePath.
		// Command.commandNamePath is filled in when the parent's Command.startApp is called
		// and startApp is only called when you actually try to run that command or one of
		// its subcommands. So we run "bytemark <command> help" on all commands that have
		// subcommands in order to get every subcommand to have a correct Command.commandPath

		if c.Subcommands != nil && len(c.Subcommands) > 0 {
			_ = cliapp.Run([]string{"bytemark.test", c.Name, "help"})
		}
	}

	buf := bytes.Buffer{}
	cliapp.Metadata["buf"] = &buf
	cliapp.Metadata["debugWriter"] = &TestWriter{t}

	cliapp.Writer = &buf

	return
}

type TestWriter struct {
	t *testing.T
}

func (tw *TestWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}

// baseTestAuthSetup sets up a 'regular' test - with auth, no yubikey.
// user is test-user
func BaseTestAuthSetup(t *testing.T, admin bool, commands []cli.Command) (config *mocks.Config, c *mocks.Client, cliapp *cli.App) {
	config, c, cliapp = BaseTestSetup(t, admin, commands)

	config.When("Get", "account").Return("test-account")
	config.When("GetIgnoreErr", "token").Return("test-token")
	config.When("GetIgnoreErr", "user").Return("test-user")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "2fa-otp").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	return
}

func traverseAllCommands(cmds []cli.Command, fn func(cli.Command)) {
	if cmds == nil {
		return
	}
	for _, c := range cmds {
		fn(c)
		traverseAllCommands(c.Subcommands, fn)
	}
}
