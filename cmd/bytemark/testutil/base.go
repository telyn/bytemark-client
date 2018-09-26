package testutil

import (
	"bytes"
	"errors"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

// GetBuf returns the captured output buffer for the given app.
// Use this in your tests to find out if your command's output is correct.
// TODO: add an example
func GetBuf(app *cli.App) (buf *bytes.Buffer, err error) {
	if bufInterface, ok := app.Metadata["buf"]; ok {
		if buf, ok = bufInterface.(*bytes.Buffer); ok {
			return
		}
		err = errors.New("couldn't recover the buffer - use BaseTestSetup to create your cli.Apps")
		return
	}
	err = errors.New("couldn't recover the buffer - use BaseTestSetup to create your cli.Apps")
	return
}

// AssertOutput fails the test unless the app under test has output exactly the
// expected string.
func AssertOutput(t *testing.T, app *cli.App, expected string) {
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

// AssertOutputMatch fails the test unless the app under test has produced
// output which matches the supplied regex.
func AssertOutputMatch(t *testing.T, app *cli.App, expected regexp.Regexp) {
	buf, err := GetBuf(app)

	if err != nil {
		t.Errorf("Couldn't recover the buffer - use BaseTestSetup to create your cli.Apps")
		return
	}
	if !expected.Match(buf.Bytes()) {
		return
	}
}

// BaseTestSetup constructs mock config and client and produces a cli.App with the given commands.
func BaseTestSetup(t *testing.T, admin bool, commands []cli.Command) (conf *mocks.Config, client *mocks.Client, cliapp *cli.App) {
	cli.OsExiter = func(_ int) {}
	conf = new(mocks.Config)
	client = new(mocks.Client)
	conf.When("GetBool", "admin").Return(admin, nil)
	conf.When("GetV", "output-format").Return(config.Var{"output-format", "human", "CODE"})

	cliapp, err := app.BaseAppSetup(app.GlobalFlags(), commands)
	if err != nil {
		t.Fatal(err)
	}
	app.SetClientAndConfig(cliapp, client, conf)

	buf := bytes.Buffer{}
	cliapp.Metadata["buf"] = &buf
	cliapp.Metadata["debugWriter"] = &TestWriter{t}

	fixCommandFullName(cliapp, commands)

	cliapp.Writer = &buf
	log.Writer = &buf
	log.ErrWriter = &buf

	return
}

// fixCommandFullName ensures that Command.FullName works for all commands in the slice.
// see the comment inside for the reasoning behind it
func fixCommandFullName(cliapp *cli.App, commands []cli.Command) {
	// discard output during this setup
	oldWriter := cliapp.Writer
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
	cliapp.Writer = oldWriter
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
	config.When("GetIgnoreErr", "impersonate").Return("")

	c.When("GetSessionFactors").Return([]string{"username", "password"})
	c.When("GetSessionUser").Return("test-user")
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	return
}

// TraverseAllCommands goes through all the commands it is supplied.
func TraverseAllCommands(cmds []cli.Command, fn func(cli.Command)) {
	if cmds == nil {
		return
	}
	for _, c := range cmds {
		fn(c)
		TraverseAllCommands(c.Subcommands, fn)
	}
}

// TraverseAllCommandsWithContext adds a more details such as the parent command to commands so we can find the offender easier.
func TraverseAllCommandsWithContext(cmds []cli.Command, name string, fn func(fullCommandString string, command cli.Command)) {
	if cmds == nil {
		return
	}
	for _, c := range cmds {
		subName := name + " " + c.FullName()
		fn(subName, c)
		TraverseAllCommandsWithContext(c.Subcommands, subName, fn)
	}
}
