package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
)

func TestResetCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("ResetVirtualMachine", vmn).Times(1)

	err := app.Run(strings.Split("bytemark reset test-server.test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestRestartCommandTable(t *testing.T) {
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
			input: "test-server --appliance",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "default",
				Account:        "default-account",
			},
			shouldErr: true,
		}, {
			name:  "RestartWithApplianceFlag",
			input: "test-server --appliance rescue",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "default",
				Account:        "default-account",
			},
			applianceBoot: true,
		}, {
			name:  "RestartWithRescueFlag",
			input: "test-server --rescue",
			vmname: lib.VirtualMachineName{
				VirtualMachine: "test-server",
				Group:          "default",
				Account:        "default-account",
			},
			applianceBoot: true,
		}, {
			name:  "RestartWithApplianceAndRescueFlag",
			input: "test-server --appliance a --rescue",
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
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands)
			config.When("GetVirtualMachine").Return(defVM)

			client.When("ShutdownVirtualMachine", test.vmname, true).Times(1)
			client.When("GetVirtualMachine", test.vmname).Return(brain.VirtualMachine{PowerOn: false})

			if test.applianceBoot == true {
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

			args := fmt.Sprintf("bytemark restart %s", test.input)
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

// func TestRestartCommand(t *testing.T) {
// 	is := is.New(t)
// 	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

// 	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

// 	config.When("GetVirtualMachine").Return(defVM)

// 	c.When("ShutdownVirtualMachine", vmn, true).Times(1)
// 	c.When("GetVirtualMachine", vmn).Return(brain.VirtualMachine{PowerOn: false})
// 	c.When("StartVirtualMachine", vmn).Times(1)

// 	err := app.Run(strings.Split("bytemark restart test-server.test-group.test-account", " "))
// 	is.Nil(err)
// 	if ok, err := c.Verify(); !ok {
// 		t.Fatal(err)
// 	}
// }
func TestShutdownCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("ShutdownVirtualMachine", vmn, true).Times(1)
	c.When("GetVirtualMachine", vmn).Return(brain.VirtualMachine{PowerOn: false})

	err := app.Run(strings.Split("bytemark shutdown test-server.test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStartCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)
	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("StartVirtualMachine", vmn).Times(1)

	err := app.Run(strings.Split("bytemark start test-server.test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
func TestStopCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	vmn := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}

	config.When("GetVirtualMachine").Return(defVM)

	c.When("StopVirtualMachine", vmn).Times(1)

	err := app.Run(strings.Split("bytemark stop test-server.test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
