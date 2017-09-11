package main

import (
	"bytes"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	mock "github.com/maraino/go-mock"
	"github.com/urfave/cli"
)

var defVM = lib.VirtualMachineName{Group: "default", Account: "default-account"}
var defGroup = lib.GroupName{Group: "default", Account: "default-account"}

func init() {
	// If we are testing, we want to override the OsExiter,
	// so we can actually test errors returned from actions
	cli.OsExiter = func(c int) {}
}

func baseTestSetup(t *testing.T, admin bool) (config *mocks.Config, client *mocks.Client, app *cli.App) {
	config = new(mocks.Config)
	client = new(mocks.Client)
	config.When("GetBool", "admin").Return(admin, nil)
	config.When("GetV", "output-format").Return(util.ConfigVar{"output-format", "human", "CODE"})

	app, err := baseAppSetup(globalFlags(), config)
	if err != nil {
		t.Fatal(err)
	}
	app.Metadata = map[string]interface{}{
		"client": client,
		"config": config,
	}

	app.Writer = ioutil.Discard
	for _, c := range commands {
		//config.When("Get", "token").Return("no-not-a-token")

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
func baseTestAuthSetup(t *testing.T, admin bool) (config *mocks.Config, c *mocks.Client, app *cli.App) {
	config, c, app = baseTestSetup(t, admin)

	config.When("Get", "account").Return("test-account")
	config.When("GetIgnoreErr", "token").Return("test-token")
	config.When("GetIgnoreErr", "user").Return("test-user")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetIgnoreErr", "2fa-otp").Return("")

	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	return config, c, app
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
		IPRanges:  make([]brain.IPRange, 0),
	}
}

func getFixtureIPRange() brain.IPRange {
	return brain.IPRange{
		ID:        1,
		Spec:      "192.168.1.1/28",
		VLANNum:   1,
		Zones:     make([]string, 0),
		Available: big.NewInt(11),
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

func getFixtureDisc() brain.Disc {
	return brain.Disc{
		ID:    132,
		Label: "disc.sata-1.132",
	}
}

func getFixtureGroup() brain.Group {
	vms := make([]brain.VirtualMachine, 1)
	vms[0] = getFixtureVM()

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

func resetOneMockedFunction(functions []*mock.MockFunction, name string, args ...interface{}) (out []*mock.MockFunction) {
	out = make([]*mock.MockFunction, 0, len(functions))
	//fmt.Printf("Time to reset one mocked function - %q\n", name)
	for _, f := range functions {
		if f.Name != name {
			//fmt.Printf("%q - carrying on\n", f.Name)
			out = append(out, f)
			continue
		}
		//fmt.Printf("Same name - %s\n", name)
		if len(args) == len(f.Arguments) {
			//fmt.Println("args same len")
			foundDifferent := false
			for i := range args {
				if args[i] != f.Arguments[i] {
					//fmt.Printf("arg %d expecting %v got %v\n", i, args[i], f.Arguments[i])
					foundDifferent = true
					break
				}
			}
			if !foundDifferent {
				//fmt.Println("args same, removing this fn")
			} else {
				out = append(out, f)
			}
		} else {
			//fmt.Println("different arg lens - appending")
			out = append(out, f)
		}
	}
	return
}
