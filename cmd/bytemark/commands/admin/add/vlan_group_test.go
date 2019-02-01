package add_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
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
			input: "add vlan group test-group.test-account",
		}, {
			name:  "alias no num",
			input: "add vlan-group test-group.test-account",
		}, {
			name:            "no num",
			input:           "add vlan group test-group.test-account 19",
			expectedVLANNum: 19,
		}, {
			name:            "alias no num",
			input:           "add vlan-group test-group.test-account 19",
			expectedVLANNum: 19,
		}, {
			name:  "err",
			input: "add vlan-group test-group.test-account",
			err:   errors.New("group name already used"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

			config.When("GetGroup").Return(testutil.DefGroup).Times(1)

			group := pathers.GroupName{
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
