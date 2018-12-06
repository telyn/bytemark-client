package app

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func baseTestSetup(t *testing.T, admin bool, commands []cli.Command) (conf *mocks.Config, client *mocks.Client, app *cli.App) {
	conf = new(mocks.Config)
	client = new(mocks.Client)
	conf.When("GetBool", "admin").Return(admin, nil)
	conf.When("GetV", "output-format").Return(config.Var{Name: "output-format", Value: "human", Source: "CODE"})

	app, err := BaseAppSetup(GlobalFlags(), commands)
	if err != nil {
		t.Fatal(err)
	}
	app.Metadata = map[string]interface{}{
		"client": client,
		"config": conf,
	}

	app.Writer = ioutil.Discard
	for _, c := range commands {
		//conf.When("Get", "token").Return("no-not-a-token")

		// the issue is that Command.FullName() is dependent on Command.commandNamePath.
		// Command.commandNamePath is filled in when the parent's Command.startApp is called
		// and startApp is only called when you actually try to run that command or one of
		// its subcommands. So we run "bytemark <command> help" on all commands that have
		// subcommands in order to get every subcommand to have a correct Command.commandPath

		if c.Subcommands != nil && len(c.Subcommands) > 0 {
			_ = app.Run([]string{"bytemark.test", c.Name, "help"})
		}
	}

	buf := bytes.Buffer{}
	app.Metadata["buf"] = &buf
	app.Metadata["debugWriter"] = &TestWriter{t}

	app.Writer = &buf

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
func BaseTestAuthSetup(t *testing.T, admin bool, commands []cli.Command) (conf *mocks.Config, c *mocks.Client, app *cli.App) {
	conf, c, app = baseTestSetup(t, admin, commands)

	conf.When("Get", "account").Return("test-account")
	conf.When("GetIgnoreErr", "token").Return("test-token")
	conf.When("GetIgnoreErr", "user").Return("test-user")
	conf.When("GetIgnoreErr", "yubikey").Return("")
	conf.When("GetIgnoreErr", "2fa-otp").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	return conf, c, app
}
