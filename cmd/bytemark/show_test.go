package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
)

type showAccountTest struct {
	// the command line to call
	Input string
	// the account name that's in the config, either because it's what was set as a global flag or it's in the config dir
	ConfigAccount string
	// the account that should be attempted to be got - if blank, expect GetDefaultAccount
	AccountToGet string
	ShouldErr    bool
}

// mkExpectedOutput loops over every value in template, running fmt.Sprintf to replace all %s with substitute
// in this way it can be used to prepare an expected output array. See showTestAccountOutput and TestShowAccount for some overly complicated example usage.
func mkExpectedOutput(template map[string]string, substitute string) map[string]string {
	m := make(map[string]string)
	for k, v := range template {
		toReplace := strings.Count(v, "%s")
		fmtArgs := make([]interface{}, toReplace)
		for i := 0; i < toReplace; i++ {
			fmtArgs[i] = substitute
		}
		m[k] = fmt.Sprintf(v, fmtArgs...)
	}
	return m
}

var showAccountTestOutput = map[string]string{
	"table": "+-----------+------+-----------+--------+\n| BillingID | Name | Suspended | Groups |\n+-----------+------+-----------+--------+\n|       213 | %s | false     |        |\n+-----------+------+-----------+--------+\n",
	"json":  "{\n    \"name\": \"%s\",\n    \"owner\": {\n        \"username\": \"\",\n        \"email\": \"\",\n        \"password\": \"\",\n        \"firstname\": \"%s\",\n        \"surname\": \"Testo\",\n        \"address\": \"\",\n        \"city\": \"\",\n        \"postcode\": \"\",\n        \"country\": \"\",\n        \"phone\": \"\"\n    },\n    \"technical_contact\": {\n        \"username\": \"\",\n        \"email\": \"\",\n        \"password\": \"\",\n        \"firstname\": \"%s\",\n        \"surname\": \"Testo\",\n        \"address\": \"\",\n        \"city\": \"\",\n        \"postcode\": \"\",\n        \"country\": \"\",\n        \"phone\": \"\"\n    },\n    \"billing_id\": 213,\n    \"brain_id\": 112,\n    \"card_reference\": \"\",\n    \"groups\": [],\n    \"suspended\": false\n}\n",
	"human": "213 - %s\n",
}

func TestShowAccount(t *testing.T) {
	baseShowAccountSetup := func(c *mocks.Client, config *mocks.Config, configAccount, outputFormat string) {
		config.Mock.Functions = resetOneMockedFunction(config.Mock.Functions, "GetV", "output-format")
		config.When("GetV", "output-format").Return(util.ConfigVar{"output-format", outputFormat, "FLAG output-format"})
		config.When("GetIgnoreErr", "account").Return(configAccount)
	}

	// These tests dont actually work, probably something more to do with the mocking of server than anything as you cannot get account that isnt your own
	// TODO: Fix these tests or remove the ones that are useless?
	tests := []showAccountTest{
		{ // 0
			Input:         "show account",
			ConfigAccount: "",
		},
		{ // 1
			Input:         "show account",
			ConfigAccount: "secg",
			AccountToGet:  "secg",
		},
		{ // 2
			Input:         "--account caan show account",
			ConfigAccount: "",
			AccountToGet:  "caan",
		},
		{ // 3
			Input:         "--account caan show account",
			ConfigAccount: "jast",
			AccountToGet:  "caan",
		},
		{ // 4
			Input:         "show account thay",
			ConfigAccount: "jast",
			AccountToGet:  "thay",
		},
		{ // 5
			Input:         "show account --account caan",
			ConfigAccount: "jast",
			AccountToGet:  "caan",
		},
	}

	runTest := func(i int, test showAccountTest) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("TestShowAccountCommand %d panicked: %v %s", i, err, debug.Stack())
			}
		}()

		expectedOutput := mkExpectedOutput(showAccountTestOutput, test.AccountToGet)
		if test.AccountToGet == "" {
			expectedOutput = mkExpectedOutput(showAccountTestOutput, "defa")
		}
		for _, format := range []string{"table", "json", "human"} {
			t.Logf("show account %d %s", i, format)

			config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
			baseShowAccountSetup(c, config, test.ConfigAccount, format)
			if test.AccountToGet == "" {
				c.When("GetAccount", "").Return(lib.Account{
					Name:      "defa",
					BrainID:   112,
					BillingID: 213,
					Groups:    []brain.Group{},
					Owner: billing.Person{
						FirstName: "defa",
						LastName:  "Testo",
					},
					TechnicalContact: billing.Person{
						FirstName: "defa",
						LastName:  "Testo",
					},
				}).Times(1)
			} else {
				// wat
				c.When("GetAccount", test.AccountToGet).Return(lib.Account{
					Name:      test.AccountToGet,
					BrainID:   112,
					BillingID: 213,
					Groups:    []brain.Group{},
					Owner: billing.Person{
						FirstName: test.AccountToGet,
						LastName:  "Testo",
					},
					TechnicalContact: billing.Person{
						FirstName: test.AccountToGet,
						LastName:  "Testo",
					},
				}).Times(1)
			}

			args := fmt.Sprintf("bytemark --output-format=%s %s", format, test.Input)
			//t.Logf("TestShowAccountCommand %d args: %s", i, args)
			err := app.Run(strings.Split(args, " "))
			if !test.ShouldErr && err != nil {
				t.Errorf("TestShowAccountCommand %d shouldn't err, but did: %T{%s}", i, err, err.Error())
			} else if test.ShouldErr && err == nil {
				t.Errorf("TestShowAccountCommand %d should err, but didn't", i)
			}
			if ok, vErr := c.Verify(); !ok {
				t.Errorf("TestShowAccountCommand %d client failed to verify: %s", i, vErr.Error())
			}
			if expected, ok := expectedOutput[format]; ok {
				testutil.AssertOutput(t, fmt.Sprintf("TestShowAccountCommand %d (%s)", i, format), app, expected)
			} else {
				buf, err := testutil.GetBuf(app)
				if err != nil {
					t.Errorf("TestShowAccountCommand %d didn't have an expected %s output. Also %s", i, format, err.Error())
				} else {
					t.Errorf("TestShowAccountCommand %d didn't have an expected %s output. Maybe it should be %q", i, format, buf.String())
				}
			}
		}
	}

	for i, test := range tests {
		runTest(i, test)
	}

}

func TestShowGroupCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	config.When("GetGroup").Return(defGroup)
	gpname := lib.GroupName{Group: "test-group", Account: "test-account"}

	group := getFixtureGroup()
	c.When("GetGroup", gpname).Return(&group, nil).Times(1)

	err := app.Run(strings.Split("bytemark show group test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowServerCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	config.When("GetVirtualMachine").Return(defVM)
	vmname := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}
	vm := getFixtureVM()
	c.When("GetVirtualMachine", vmname).Return(&vm, nil).Times(1)

	err := app.Run(strings.Split("bytemark show server test-server.test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowPrivileges(t *testing.T) {

	tests := []struct {
		privs     brain.Privileges
		args      string
		user      string
		shouldErr bool
	}{
		{
			privs:     brain.Privileges{{VirtualMachineID: 643, Username: "satan", Level: brain.VMConsolePrivilege}},
			user:      "",
			args:      "bytemark --admin show privileges",
			shouldErr: false,
		},
	}

	for i, test := range tests {
		_, c, app := testutil.BaseTestAuthSetup(t, true, commands)
		c.When("GetPrivileges", test.user).Return(test.privs, nil)

		err := app.Run(strings.Split(test.args, " "))
		if test.shouldErr && err == nil {
			t.Errorf("TestShowPrivilege %d should err and didn't", i)
		} else if !test.shouldErr && err != nil {
			t.Errorf("TestShowPrivilege %d shouldn't err, but: %s", i, err)
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}

func TestShowDisc(t *testing.T) {
	tests := []struct {
		disc          brain.Disc
		err           error
		shouldErr     bool
		args          string
		vm            lib.VirtualMachineName
		discLabelOrID string
	}{
		{ // 0 - fully specified flags
			disc: brain.Disc{
				ID:           42,
				Label:        "faff",
				StorageGrade: "sata",
			},
			args:          "bytemark show disc --server fliff --disc faff",
			vm:            lib.VirtualMachineName{VirtualMachine: "fliff", Group: "default", Account: "default-account"},
			discLabelOrID: "faff",
		}, { // 1 - flags from args
			disc: brain.Disc{
				ID:           42,
				Label:        "faff",
				StorageGrade: "sata",
			},
			args:          "bytemark show disc fliff faff",
			vm:            lib.VirtualMachineName{VirtualMachine: "fliff", Group: "default", Account: "default-account"},
			discLabelOrID: "faff",
		}, { // 2 - no disc
			args:      "bytemark show disc fliff",
			shouldErr: true,
		}, { // 3 --server but no --disc
			args:      "bytemark show disc --server fliff",
			shouldErr: true,
		}, { // 4 - numeric --disc with no server
			disc: brain.Disc{
				ID:           42,
				Label:        "faff",
				StorageGrade: "sata",
			},
			vm:            lib.VirtualMachineName{},
			args:          "bytemark show disc --disc 42",
			discLabelOrID: "42",
		},
	}
	var i = 0
	var test = tests[0]

	defer func() {
		if msg := recover(); msg != nil {
			t.Errorf("TestShowDisc %d panic: %v", i, msg)
		}
	}()

	for i, test = range tests {
		config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
		if test.shouldErr && test.err == nil {
			config, c, app = testutil.BaseTestSetup(t, false, commands)
		}
		config.When("GetVirtualMachine").Return(defVM)
		c.When("GetDisc", test.vm, test.discLabelOrID).Return(test.disc, test.err)

		err := app.Run(strings.Split(test.args, " "))
		if test.shouldErr != (err != nil) {
			nt := ""
			if test.shouldErr {
				nt = "n't"
			}
			t.Errorf("%s should%s error, got %s", testutil.Name(i), nt, err)
		}
		if ok, vErr := c.Verify(); !ok {
			t.Errorf("%s client failed to verify: %s", testutil.Name(i), vErr.Error())
		}
	}
}
