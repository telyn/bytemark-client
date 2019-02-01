package delete_test

import (
	"strings"
	"testing"

	appPkg "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	mock "github.com/maraino/go-mock"
)

func TestDeleteGroup(t *testing.T) {
	// ok i need to remember how to test prompts.
	tests := []struct {
		name string

		serverInGroup  bool
		command        string
		shouldPrompt   bool
		promptResponse string
		shouldCall     bool
		shouldErr      bool
	}{{
		name:           "empty group no recurse no force Y",
		command:        "delete group test-group",
		shouldPrompt:   true,
		promptResponse: "y",
		shouldCall:     true,
	}, {
		name:           "empty group no recurse no force N",
		command:        "delete group test-group",
		shouldPrompt:   true,
		promptResponse: "n",
		shouldCall:     false,
		shouldErr:      true,
	}, {
		name:         "empty group no recurse +force",
		command:      "delete group --force test-group",
		shouldPrompt: false,
		shouldCall:   true,
	}, {
		name:           "empty group +recurse no force Y",
		command:        "delete group --recursive test-group",
		shouldPrompt:   true,
		promptResponse: "y",
		shouldCall:     true,
	}, {
		name:           "empty group +recurse no force N",
		command:        "delete group --recursive test-group",
		shouldPrompt:   true,
		promptResponse: "n",
		shouldCall:     false,
		shouldErr:      true,
	}, {
		name:         "empty group +recurse force",
		command:      "delete group --recursive --force test-group",
		shouldPrompt: false,
		shouldCall:   true,
	}, {
		name:          "server group no recurse no force",
		command:       "delete group test-group",
		serverInGroup: true,
		shouldPrompt:  false,
		shouldErr:     true,
		shouldCall:    false,
	}, {
		name:          "server group no recurse +force",
		command:       "delete group --force test-group",
		serverInGroup: true,
		shouldPrompt:  false,
		shouldErr:     true,
		shouldCall:    false,
	}, {
		name:           "server group +recurse no force Y",
		command:        "delete group --recursive test-group",
		serverInGroup:  true,
		shouldPrompt:   true,
		promptResponse: "y",
		shouldCall:     true,
	}, {
		name:           "server group +recurse no force N",
		command:        "delete group --recursive test-group",
		serverInGroup:  true,
		shouldPrompt:   true,
		promptResponse: "n",
		shouldCall:     false,
		shouldErr:      true,
	}, {
		name:          "server group +recurse force",
		command:       "delete group --recursive --force test-group",
		serverInGroup: true,
		shouldPrompt:  false,
		shouldCall:    true,
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testPrompter := mocks.Prompter{}
			if test.shouldPrompt {
				testPrompter.When("Prompt", mock.Any).Return(test.promptResponse)
			}

			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			appPkg.SetPrompter(app, &testPrompter)

			group := brain.Group{
				Name: "test-group",
				ID:   9000,
			}

			if test.serverInGroup {
				group.VirtualMachines = brain.VirtualMachines{
					brain.VirtualMachine{
						Name: "test-vm",
					},
				}
				client.When("DeleteVirtualMachine", lib.VirtualMachineName{
					Account:        "test-account",
					Group:          "test-group",
					VirtualMachine: "test-vm",
				}, true).Return(nil)
			}

			config.When("GetGroup").Return(pathers.GroupName{Account: "test-account"})
			client.When("GetGroup", pathers.GroupName{Group: "test-group", Account: "test-account"}).Return(group)
			if test.shouldCall {
				client.When("DeleteGroup", pathers.GroupName{Group: "test-group", Account: "test-account"}).Return(nil)
			}

			err := app.Run(strings.Split("bytemark "+test.command, " "))
			if err != nil && !test.shouldErr {
				t.Errorf("Unexpected error from app.Run: %s", err)
			} else if err == nil && test.shouldErr {
				t.Error("Expected error but did not get one :-C")
			}
			if ok, err := testPrompter.Verify(); !ok {
				t.Fatal(err)
			}
			if ok, err := client.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}
