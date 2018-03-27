package admin_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/util"
	"github.com/cheekybits/is"
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestCreateVLANGroup(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedVLANNum int
		err             error
	}{
		{
			name:  "no num",
			input: "create vlan group test-group.test-account",
		}, {
			name:  "alias no num",
			input: "create vlan-group test-group.test-account",
		}, {
			name:            "no num",
			input:           "create vlan group test-group.test-account 19",
			expectedVLANNum: 19,
		}, {
			name:            "alias no num",
			input:           "create vlan-group test-group.test-account 19",
			expectedVLANNum: 19,
		}, {
			name:  "err",
			input: "create vlan-group test-group.test-account",
			err:   errors.New("group name already used"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

			config.When("GetGroup").Return(defGroup).Times(1)

			group := lib.GroupName{
				Group:   "test-group",
				Account: "test-account",
			}
			c.When("AdminCreateGroup", group, test.expectedVLANNum).Return(test.err).Times(1)

			err := app.Run(strings.Split("bytemark "+test.input, " "))
			if test.err != nil && err == nil {
				t.Error("expected error but received none")
			} else if test.err == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if ok, err := c.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateIPRange(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateIPRange", "192.168.3.0/28", 14).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark create ip range 192.168.3.0/28 14", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateIPRangeError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateIPRange", "192.168.3.0/28", 18).Return(fmt.Errorf("Error creating IP range")).Times(1)

	err := app.Run(strings.Split("bytemark create ip range 192.168.3.0/28 18", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateUser", "uname", "cluster_su").Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark create user uname cluster_su", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateUserError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	c.When("CreateUser", "uname", "cluster_su").Return(fmt.Errorf("Error creating user")).Times(1)

	err := app.Run(strings.Split("bytemark create user uname cluster_su", " "))
	is.NotNil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateMigration(t *testing.T) {
	tests := []struct {
		name      string
		args      string
		exp       brain.MigrationJob
		shouldErr bool
	}{
		{
			name: "OneDisc",
			args: "--disc 1",
			exp: brain.MigrationJob{
				Args: brain.MigrationJobSpec{
					Sources: brain.MigrationJobLocations{
						Discs: []util.NumberOrString{"1"},
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{},
					},
					Destinations: brain.MigrationJobLocations{
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{},
					},
				},
				Queue: brain.MigrationJobQueue{
					Discs: []int{1},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(y *testing.T) {
			_, client, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)
			postReq := &mocks.Request{
				T:              t,
				StatusCode:     201,
				ResponseObject: test.exp,
			}
			client.When("BuildRequest", "POST", lib.BrainEndpoint, "/admin/migration_jobs%s", []string{""}).Return(postReq).Times(1)
			args := strings.Split("bytemark --admin create migration "+test.args, " ")
			err := app.Run(args)
			if !test.shouldErr && err != nil {
				t.Errorf("shouldn't err, but did: %T{%s}", err, err.Error())
			} else if test.shouldErr && err == nil {
				t.Fatal("should err, but didn't")
				return
			}
			if ok, err := client.Verify(); !ok {
				t.Fatal(err)
			}
			postReq.AssertRequestObjectEqual(test.exp.Args)

		})
	}
}
