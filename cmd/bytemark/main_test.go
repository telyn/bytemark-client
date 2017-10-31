package main

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
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
