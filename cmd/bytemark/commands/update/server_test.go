package update_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

func TestUpdateServer(t *testing.T) {
	defVM := lib.VirtualMachineName{
		Group:   "default",
		Account: "default-account",
	}
	testVM := lib.VirtualMachineName{
		VirtualMachine: "test",
		Group:          "default",
		Account:        "default-account",
	}
	type move struct {
		expected bool
		newName  lib.VirtualMachineName
	}
	tests := []struct {
		name          string
		args          string
		vmName        lib.VirtualMachineName
		hwProfile     string
		hwProfileLock bool
		move          move
		shouldErr     bool
	}{
		{
			name:   "RenameVM",
			args:   "--new-name new --server test",
			vmName: testVM,
			move: move{
				expected: true,
				newName: lib.VirtualMachineName{
					VirtualMachine: "new",
					Group:          "default",
					Account:        "default-account",
				},
			},
		},
		{
			name:      "ChangeHWProfile",
			args:      "--hw-profile foo --server test",
			vmName:    testVM,
			hwProfile: "foo",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			config.When("GetVirtualMachine").Return(defVM)
			config.When("Force").Return(true)

			if test.hwProfile != "" {
				client.When("SetVirtualMachineHardwareProfile", test.vmName, test.hwProfile, []bool{test.hwProfileLock}).Return(nil).Times(1)
			}

			if test.move.expected {
				client.When("MoveVirtualMachine", test.vmName, test.move.newName).Return(nil).Times(1)
			}

			args := strings.Split("bytemark update server "+test.args, " ")
			err := app.Run(args)
			if test.shouldErr && err == nil {
				t.Fatal("should error")
			} else if !test.shouldErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok, err := client.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}
