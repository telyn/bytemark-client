package update_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func TestUpdateServer(t *testing.T) {
	defVM := lib.VirtualMachineName{
		Group:   "default",
		Account: "default-account",
	}
	testVMName := lib.VirtualMachineName{
		VirtualMachine: "test",
		Group:          "default",
		Account:        "default-account",
	}
	testVM := brain.VirtualMachine{
		Name:     "test",
		Hostname: "test.default",
		GroupID:  1,
		Memory:   1024,
	}
	type move struct {
		expected bool
		newName  lib.VirtualMachineName
	}
	tests := []struct {
		name          string
		args          string
		vmName        lib.VirtualMachineName
		vm            brain.VirtualMachine
		memory        int
		hwProfile     string
		hwProfileLock bool
		move          move
		shouldErr     bool
	}{
		{
			name:   "RenameVM",
			args:   "--new-name new --server test",
			vmName: testVMName,
			vm:     testVM,
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
			name:   "MoveVMGroup",
			args:   "--new-name test.foo --server test",
			vmName: testVMName,
			vm:     testVM,
			move: move{
				expected: true,
				newName: lib.VirtualMachineName{
					VirtualMachine: "test",
					Group:          "foo",
					Account:        "default-account",
				},
			},
		},
		{
			name:   "MoveVMAccount",
			args:   "--new-name test.default.other-account --server test",
			vmName: testVMName,
			vm:     testVM,
			move: move{
				expected: true,
				newName: lib.VirtualMachineName{
					VirtualMachine: "test",
					Group:          "default",
					Account:        "other-account",
				},
			},
		},
		{
			name:      "HWProfileLockErrorsWithoutHWProfile",
			args:      "--hw-profile-lock --server test",
			vmName:    testVMName,
			vm:        testVM,
			shouldErr: true,
		},
		{
			name:      "ChangeHWProfile",
			args:      "--hw-profile foo --server test",
			vmName:    testVMName,
			vm:        testVM,
			hwProfile: "foo",
		},
		{
			name:          "ChangeHWProfileWithLock",
			args:          "--hw-profile-lock --hw-profile foo --server test",
			vmName:        testVMName,
			vm:            testVM,
			hwProfile:     "foo",
			hwProfileLock: true,
		},
		{
			name:   "ChangeMemory",
			args:   "--force --memory 10 --server test",
			vmName: testVMName,
			vm:     testVM,
			memory: 10,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			config.When("GetVirtualMachine").Return(defVM)
			client.When("GetVirtualMachine", test.vmName).Return(&test.vm).Times(1)
			config.When("Force").Return(true)

			if test.hwProfile != "" {
				client.When("SetVirtualMachineHardwareProfile", test.vmName, test.hwProfile, []bool{test.hwProfileLock}).Return(nil).Times(1)
			}

			if test.move.expected {
				client.When("MoveVirtualMachine", test.vmName, test.move.newName).Return(nil).Times(1)
			}

			if test.memory > 0 {
				client.When("SetVirtualMachineMemory", test.vmName, test.memory * 1024).Return(nil).Times(1)
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
