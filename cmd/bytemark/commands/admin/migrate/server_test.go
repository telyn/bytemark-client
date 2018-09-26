package migrate_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestMigrateServer(t *testing.T) {
	tests := []struct {
		name       string
		args       string
		head       string
		vm         lib.VirtualMachineName
		vmid       int
		migrateErr error
		shouldErr  bool
	}{{
		name: "with new head",
		head: "stg-h1",
		vm: lib.VirtualMachineName{
			VirtualMachine: "vm123",
			Group:          "group",
			Account:        "account",
		},

		args: "migrate vm vm123.group.account stg-h1",
	}, {
		name: "without new head",
		vm: lib.VirtualMachineName{
			VirtualMachine: "vm122",
			Group:          "group",
			Account:        "account",
		},
		args: "migrate vm vm122.group.account",
	}, {
		name: "error",
		vm: lib.VirtualMachineName{
			VirtualMachine: "vm121",
			Group:          "group",
			Account:        "account",
		},
		args:       "migrate vm vm121.group.account",
		migrateErr: fmt.Errorf("all the heads are down oh no"),
		shouldErr:  true,
	}}

	for _, test := range tests {
		testutil.CommandT{
			Name:      test.name,
			Auth:      true,
			Admin:     true,
			Args:      test.args,
			ShouldErr: test.shouldErr,
			Commands:  admin.Commands,
		}.Run(t, func(t *testing.T, config *mocks.Config, client *mocks.Client, app *cli.App) {
			config.When("GetVirtualMachine").Return(lib.VirtualMachineName{
				VirtualMachine: "test-vm",
				Group:          "test-group",
				Account:        "test-account",
			})
			client.When("GetVirtualMachine", test.vm).Return(brain.VirtualMachine{
				ID:       13412,
				Hostname: "real.cool.vm.uk0.bigv.io",
			})

			client.When("MigrateVirtualMachine", test.vm, test.head).Return(test.migrateErr).Times(1)
		})
	}
}
