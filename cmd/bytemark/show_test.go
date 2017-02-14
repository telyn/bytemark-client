package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestShowGroupCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestSetup(t, false)

	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetGroup").Return(&defGroup)
	gpname := lib.GroupName{Group: "test-group", Account: "test-account"}
	c.When("ParseGroupName", "test-group.test-account", []*lib.GroupName{&defGroup}).Return(&gpname)
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
