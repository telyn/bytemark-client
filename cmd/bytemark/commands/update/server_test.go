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
		Memory:   2048,
		Cores:    2,
	}
	type move struct {
		expected bool
		newName  lib.VirtualMachineName
	}
	tests := []struct {
		name      string
		args      string
		vmName    lib.VirtualMachineName
		vm        brain.VirtualMachine
		memory    int
		hwProfile string
		lock      bool
		unlock    bool
		cores     int
		cdrom     string
		eject     bool
		move      move
		shouldErr bool
	}{
		{
			name:   "renaming server",
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
			name:   "moving server across groups",
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
			name:   "moving server across accounts",
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
			name:      "--hwprofile",
			args:      "--hwprofile foo --server test",
			vmName:    testVMName,
			vm:        testVM,
			hwProfile: "foo",
		},
		{
			name:      "--lock-hwprofile with --hwprofile",
			args:      "--lock-hwprofile --hwprofile foo --server test",
			vmName:    testVMName,
			vm:        testVM,
			hwProfile: "foo",
			lock:      true,
		},
		{
			name:   "--unlock-hwprofile",
			args:   "--unlock-hwprofile --server test",
			vmName: testVMName,
			vm:     testVM,
			unlock: true,
		},
		{
			name:   "decreasing memory doesn't need --force",
			args:   "--memory 1 --server test",
			vmName: testVMName,
			vm:     testVM,
			memory: 1,
		},
		{
			name:      "increasing memory needs --force",
			args:      "--memory 10 --server test",
			vmName:    testVMName,
			vm:        testVM,
			shouldErr: true,
		},
		{
			name:   "decreasing cores doesn't need --force",
			args:   "--cores 1 --server test",
			vmName: testVMName,
			vm:     testVM,
			cores:  1,
		},
		{
			name:      "increasing cores needs --force",
			args:      "--cores 4 --server test",
			vmName:    testVMName,
			vm:        testVM,
			shouldErr: true,
		},
		{
			name:   "ejecting cd",
			args:   "--remove-cd --server test",
			vmName: testVMName,
			vm:     testVM,
			eject:  true,
		},
		{
			name:   "inserting cd",
			args:   "--cd-url https://microsoft.com/windows.iso --server test",
			vmName: testVMName,
			vm:     testVM,
			cdrom:  "https://microsoft.com/windows.iso",
		},
		{
			name:   "ejecting and updating cd",
			args:   "--remove-cd --cd-url https://microsoft.com/windows.iso --server test",
			vmName: testVMName,
			vm:     testVM,
			cdrom:  "https://microsoft.com/windows.iso",
			eject:  true,
		},
		{
			name:      "multiple combined changes",
			args:      "--new-name new --memory 1 --hwprofile foo --cores 1 --remove-cd --cd-url https://microsoft.com/windows.iso --server test",
			vmName:    testVMName,
			vm:        testVM,
			hwProfile: "foo",
			memory:    1,
			cores:     1,
			cdrom:     "https://microsoft.com/windows.iso",
			eject:     true,
			move: move{
				expected: true,
				newName: lib.VirtualMachineName{
					VirtualMachine: "new",
					Group:          "default",
					Account:        "default-account",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("%v", err)
				}
			}()
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			config.When("GetVirtualMachine").Return(defVM)
			client.When("GetVirtualMachine", test.vmName).Return(test.vm).Times(1)

			if test.hwProfile != "" {
				client.When("SetVirtualMachineHardwareProfile", test.vmName, test.hwProfile, nil).Return(nil).Times(1)
			}
			if test.lock {
				client.When("SetVirtualMachineHardwareProfileLock", test.vmName, true).Return(nil).Times(1)
			}

			if test.unlock {
				client.When("SetVirtualMachineHardwareProfileLock", test.vmName, false).Return(nil).Times(1)
			}

			if test.move.expected {
				client.When("MoveVirtualMachine", test.vmName, test.move.newName).Return(nil).Times(1)
			}

			if test.memory > 0 {
				client.When("SetVirtualMachineMemory", test.vmName, test.memory*1024).Return(nil).Times(1)
			}

			if test.cores > 0 {
				client.When("SetVirtualMachineCores", test.vmName, test.cores).Return(nil).Times(1)
			}

			if test.eject || test.cdrom != "" {
				client.When("SetVirtualMachineCDROM", test.vmName, test.cdrom).Return(nil).Times(1)
			}

			args := strings.Split("bytemark update server "+test.args, " ")
			err := app.Run(args)
			if test.shouldErr && err == nil {
				t.Error("should error")
			} else if !test.shouldErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if ok, err := client.Verify(); !ok {
				t.Error(err)
			}
		})
	}
}
