package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
	"runtime/debug"
	"strings"
	"testing"
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

func TestShowAccountCommand(t *testing.T) {
	baseShowAccountSetup := func(c *mocks.Client, config *mocks.Config, configAccount string) {
		config.When("Get", "token").Return("test-token")
		config.When("GetIgnoreErr", "account").Return(configAccount)
		config.When("GetIgnoreErr", "yubikey").Return("")
		c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	}
	tests := []showAccountTest{
		{ // 0
			Input:         "bytemark show account",
			ConfigAccount: "",
		},
		{ // 1
			Input:         "bytemark show account",
			ConfigAccount: "sec",
			AccountToGet:  "sec",
		},
		{ // 2
			Input:         "bytemark --account caan show account",
			ConfigAccount: "",
			AccountToGet:  "caan",
		},
		{ // 3
			Input:         "bytemark --account caan show account",
			ConfigAccount: "jast",
			AccountToGet:  "caan",
		},
		{ // 4
			Input:         "bytemark show account thay",
			ConfigAccount: "jast",
			AccountToGet:  "thay",
		},
		{ // 5
			Input:         "bytemark show account --account jast",
			ConfigAccount: "jast",
			AccountToGet:  "jast",
		},
	}

	config, c := baseTestSetup(t, false)

	runTest := func(i int, test showAccountTest) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("TestShowAccountCommand %d panicked: %v %s", i, err, debug.Stack())
			}
		}()
		baseShowAccountSetup(c, config, test.ConfigAccount)
		if test.AccountToGet == "" {
			c.When("GetAccount", "").Return(&lib.Account{
				Name:      "default-account",
				BrainID:   112,
				BillingID: 213,
				Groups:    []*brain.Group{},
				Owner: &billing.Person{
					FirstName: "Doctor",
					LastName:  "Testo",
				},
				TechnicalContact: &billing.Person{
					FirstName: "Doctor",
					LastName:  "Testo",
				},
			}).Times(1)
		} else {
			c.When("GetAccount", test.AccountToGet).Return(&lib.Account{
				Name:      "default-account",
				BrainID:   112,
				BillingID: 213,
				Groups:    []*brain.Group{},
				Owner: &billing.Person{
					FirstName: "Doctor",
					LastName:  "Testo",
				},
				TechnicalContact: &billing.Person{
					FirstName: "Doctor",
					LastName:  "Testo",
				},
			}).Times(1)
		}
		err := global.App.Run(strings.Split(test.Input, " "))
		if !test.ShouldErr && err != nil {
			t.Errorf("TestShowAccountCommand %d shouldn't err, but did: %T{%s}", i, err, err.Error())
		} else if test.ShouldErr && err == nil {
			t.Errorf("TestShowAccountCommand %d should err, but didn't", i)
		}
		if ok, vErr := c.Verify(); !ok {
			t.Errorf("TestShowAccountCommand %d client failed to verify: %s", i, vErr.Error())
		}
		c.Reset()
		config.Reset()
	}

	for i, test := range tests {
		runTest(i, test)
	}

}

func TestShowGroupCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetGroup").Return(&defGroup)
	gpname := lib.GroupName{Group: "test-group", Account: "test-account"}
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	group := getFixtureGroup()
	c.When("GetGroup", &gpname).Return(&group, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark show group test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowServerCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)
	vmname := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)
	vm := getFixtureVM()
	c.When("GetVirtualMachine", &vmname).Return(&vm, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark show server test-server.test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

// TODO(telyn): show account? show user?
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
		_, c := baseTestAuthSetup(t, true)
		c.When("GetPrivileges", test.user).Return(test.privs, nil)

		err := global.App.Run(strings.Split(test.args, " "))
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
