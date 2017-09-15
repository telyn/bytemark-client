package testutil

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func GetBuf(app *cli.App) (buf *bytes.Buffer, err error) {
	if bufInterface, ok := app.Metadata["buf"]; ok {
		if buf, ok = bufInterface.(*bytes.Buffer); ok {
			return
		}
		err = errors.New("Couldn't recover the buffer - use BaseTestSetup to create your cli.Apps!")
		return
	}
	err = errors.New("Couldn't recover the buffer - use BaseTestSetup to create your cli.Apps!")
	return
}

// Asserts that the app under test has output the expected string, failing the
// test if not. identifier is a string used to identify the test - usually
// the function name of the test, plus some integer for the index of the test
// in the test-table. See show_test.go's TestShowAccountCommand for example usage
func AssertOutput(t *testing.T, identifier string, app *cli.App, expected string) {
	buf, err := GetBuf(app)

	if err == nil {
		actual := buf.String()
		if actual != expected {
			t.Errorf("expected %q, got %q", expected, actual)
		}
	} else {
		t.Error("Couldn't recover the buffer - use BaseTestSetup to create your cli.Apps!")
	}
}

// BaseTestSetup constructs mock config and client and produces a cli.App with the given commands.
func BaseTestSetup(t *testing.T, admin bool, commands []cli.Command) (config *mocks.Config, client *mocks.Client, cliapp *cli.App) {
	config = new(mocks.Config)
	client = new(mocks.Client)
	config.When("GetBool", "admin").Return(admin, nil)
	config.When("GetV", "output-format").Return(util.ConfigVar{"output-format", "human", "CODE"})

	cliapp, err := app.BaseAppSetup(app.GlobalFlags(), commands)
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

// TestWriter is a writer which writes to the test log.
// This ruins formatting on e.g. text/template renders, but at least it forces all the output to be in order
type TestWriter struct {
	t *testing.T
}

// Write pushes out the bytes as a string to the testing.T instance stored by the TestWriter using t.Log
func (tw *TestWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}

// BaseTestAuthSetup sets up a 'regular' test - with auth, no yubikey.
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
