package commands_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestRestartServerCommand(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		vmname        lib.VirtualMachineName
		shouldErr     bool
		applianceBoot bool
	}{
		{
			name:  "RestartWithoutAppliance",
			input: "test-server.test-group.test-account",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "test-group",
				Account:        "test-account",
			},
		}, {
			name:  "RestartWithoutApplianceWithDefaultAccount",
			input: "test-server.test-group",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "test-group",
				Account:        "default-account",
			},
		}, {
			name:  "RestartWithoutApplianceWithDefaultGroup",
			input: "test-server",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "default",
				Account:        "default-account",
			},
		}, {
			name:  "RestartWithApplianceFlagWithoutAppliance",
			input: "--appliance test-server",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "default",
				Account:        "default-account",
			},
			shouldErr: true,
		}, {
			name:  "RestartWithApplianceFlag",
			input: "--appliance rescue test-server",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "default",
				Account:        "default-account",
			},
			applianceBoot: true,
		}, {
			name:  "RestartWithRescueFlag",
			input: "--rescue test-server",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "default",
				Account:        "default-account",
			},
			applianceBoot: true,
		}, {
			name:  "RestartWithApplianceAndRescueFlag",
			input: "--appliance a --rescue test-server",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "default",
				Account:        "default-account",
			},
			shouldErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			config.When("GetVirtualMachine").Return(lib.VirtualMachineName{Group: "default", Account: "default-account"})
			config.When("PanelURL").Return("something.com")

			client.When("ShutdownVirtualMachine", test.vmname, true).Times(1)
			client.When("GetVirtualMachine", test.vmname).Return(brain.VirtualMachine{PowerOn: false})

			if test.applianceBoot {
				client.When("BuildRequest", "PUT", lib.Endpoint(1),
					"/accounts/%s/groups/%s/virtual_machines/%s",
					[]string{
						"default-account",
						"default",
						"test-server"}).Return(&mocks.Request{
					T:          t,
					StatusCode: 200,
				}).Times(1)
			} else {
				client.When("StartVirtualMachine", test.vmname).Times(1)
			}

			args := fmt.Sprintf("bytemark restart server %s", test.input)
			err := app.Run(strings.Split(args, " "))
			if !test.shouldErr && err != nil {
				t.Errorf("shouldn't err, but did: %T{%s}", err, err.Error())
			} else if test.shouldErr && err == nil {
				t.Errorf("should err, but didn't")
			}
			if !test.shouldErr {
				if ok, err := client.Verify(); !ok {
					t.Fatal(err)
				}
			}
		})
	}
}
