package main

import (
	"fmt"
	auth3 "github.com/BytemarkHosting/auth-client"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	mock "github.com/maraino/go-mock"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"testing"
)

var defVM = lib.VirtualMachineName{Group: "default", Account: "default-account"}
var defGroup = lib.GroupName{Group: "default", Account: "default-account"}

func TestEnsureAuth(t *testing.T) {
	tt := []struct {
		InputUsername             string
		InputPassword             string
		Input2FA                  string
		AuthWithCredentialsErrors []error
		Factors                   []string
		ExpectedError             bool
	}{
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			AuthWithCredentialsErrors: []error{nil},
			ExpectedError:             false,
		},
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			AuthWithCredentialsErrors: []error{fmt.Errorf("{}")},
			ExpectedError:             true,
		},
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			Input2FA:                  "123456",
			AuthWithCredentialsErrors: []error{fmt.Errorf("Missing 2FA"), nil}, // 2nd error as nil tests success with 2FA login
			Factors:                   []string{"2fa"},
			ExpectedError:             false,
		},
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			Input2FA:                  "123456",
			AuthWithCredentialsErrors: []error{fmt.Errorf("Missing 2FA"), fmt.Errorf("Invalid token")}, // 2nd error tests failure with 2FA token
			ExpectedError:             true,
		},
		{
			InputUsername:             "input-user",
			InputPassword:             "input-pass",
			Input2FA:                  "123456",
			AuthWithCredentialsErrors: []error{fmt.Errorf("Missing 2FA"), nil}, // 2nd error as nil means success with 2FA token
			Factors:                   []string{"missing-2fa-factor"},
			ExpectedError:             true,
		},
	}

	for _, test := range tt {
		_, c := baseTestSetup(t, false)

		configDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Errorf("Unexpected error when setting up config temp directory: %v", err)
		}
		defer func() {
			removeErr := os.RemoveAll(configDir)
			if removeErr != nil {
				t.Errorf("Could not clean up config dir: %v", removeErr)
			}
		}()

		config, err := util.NewConfig(configDir)
		if err != nil {
			t.Errorf("Unexpected error when setting up config temp directory: %v", err)
		}

		global.Config = config

		// Pretending the input comes from terminal
		global.Config.Set("user", test.InputUsername, "INTERACTION")
		global.Config.Set("pass", test.InputPassword, "TESTING")
		global.Config.Set("2fa-otp", test.Input2FA, "TESTING")

		c.When("AuthWithToken", "").Return(fmt.Errorf("Not logged in")).Times(1)

		credentials := auth3.Credentials{
			"username": test.InputUsername,
			"password": test.InputPassword,
		}

		c.When("AuthWithCredentials", credentials).Return(test.AuthWithCredentialsErrors[0]).Times(1)

		// We are supplying a 2FA token, so we want to test that flow
		if test.Input2FA != "" {
			credentials := auth3.Credentials{
				"username": test.InputUsername,
				"password": test.InputPassword,
				"2fa":      test.Input2FA,
			}
			c.When("AuthWithCredentials", credentials).Return(test.AuthWithCredentialsErrors[1]).Times(1) // Returns nil means success
		}

		// Only called if the login succeeded, so always return a token
		c.When("GetSessionToken").Return("test-token")

		c.When("GetSessionFactors").Return(test.Factors)

		err = EnsureAuth()
		if test.ExpectedError && err == nil {
			t.Error("Expecting EnsureAuth to error, but it didn't")
		} else if !test.ExpectedError && err != nil {
			t.Errorf("Not expecting EnsureAuth to error, but got %v", err)
		}

		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}

func baseTestSetup(t *testing.T, admin bool) (config *mocks.Config, client *mocks.Client) {
	config = new(mocks.Config)
	client = new(mocks.Client)
	config.When("GetBool", "admin").Return(admin, nil)
	global.Client = client
	global.Config = config

	app, err := baseAppSetup(globalFlags())
	if err != nil {
		t.Fatal(err)
	}
	global.App = app
	oldWriter := global.App.Writer
	global.App.Writer = ioutil.Discard
	for _, c := range commands {
		//config.When("Get", "token").Return("no-not-a-token")

		// the issue is that Command.FullName() is dependent on Command.commandNamePath.
		// Command.commandNamePath is filled in when the parent's Command.startApp is called
		// and startApp is only called when you actually try to run that command or one of
		// its subcommands. So we run "bytemark <command> help" on all commands that have
		// subcommands in order to get every subcommand to have a correct Command.commandPath

		if c.Subcommands != nil && len(c.Subcommands) > 0 {
			_ = global.App.Run([]string{"bytemark.test", c.Name, "help"})
		}
	}
	global.App.Writer = oldWriter
	return
}

// baseTestAuthSetup sets up a 'regular' test - with auth, no yubikey.
// user is test-user
func baseTestAuthSetup(t *testing.T, admin bool) (config *mocks.Config, c *mocks.Client) {
	config, c = baseTestSetup(t, admin)

	config.When("Get", "token").Return("test-token")
	config.When("Get", "account").Return("test-account")
	config.When("GetIgnoreErr", "user").Return("test-user")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "2fa-otp").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	return config, c
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

func getFixtureVM() brain.VirtualMachine {
	return brain.VirtualMachine{
		Name:     "test-server",
		Hostname: "test-server.test-group",
		GroupID:  1,
	}
}

func getFixtureVLAN() brain.VLAN {
	return brain.VLAN{
		ID:        1,
		Num:       1,
		UsageType: "",
		IPRanges:  make([]*brain.IPRange, 0),
	}
}

func getFixtureIPRange() brain.IPRange {
	return brain.IPRange{
		ID:        1,
		Spec:      "192.168.1.1/28",
		VLANNum:   1,
		Zones:     make([]string, 0),
		Available: 11,
	}
}

func getFixtureHead() brain.Head {
	return brain.Head{
		ID:    1,
		Label: "h1",
	}
}

func getFixtureTail() brain.Tail {
	return brain.Tail{
		ID:    1,
		Label: "t1",
	}
}

func getFixtureStoragePool() brain.StoragePool {
	return brain.StoragePool{
		Name:  "sata1",
		Label: "t1-sata1",
	}
}

func getFixtureGroup() brain.Group {
	vms := make([]*brain.VirtualMachine, 1)
	vm := getFixtureVM()
	vms[0] = &vm

	return brain.Group{
		Name:            "test-group",
		VirtualMachines: vms,
	}
}

func assertEqual(t *testing.T, test string, testNum int, name string, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("%s %d: wrong %s: expected %#v, got %#v", test, testNum, name, expected, actual)
	}
}

func checkErr(t *testing.T, name string, testNum int, shouldErr bool, err error) {
	if err == nil && shouldErr {
		t.Errorf("%s %d should error but didn't.", name, testNum)
	} else if err != nil && !shouldErr {
		t.Errorf("%s %d returned unexpected error: %s", name, testNum, err.Error())
	}
}

type Verifyer interface {
	Verify() (bool, error)
	Reset() *mock.Mock
}

func verifyAndReset(t *testing.T, name string, testNum int, v Verifyer) {
	if ok, err := v.Verify(); !ok {
		t.Errorf("%s %d %v", name, testNum, err)
	}
	v.Reset()
}
