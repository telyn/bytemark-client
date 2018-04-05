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
		name          string
		args          string
		vmName        lib.VirtualMachineName
		vm            brain.VirtualMachine
		memory        int
		hwProfile     string
		hwProfileLock bool
		cores         int
		cdrom         string
		eject         bool
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
			name:   "ChangeMemoryDoesNotNeedForceToDecrease",
			args:   "--memory 1 --server test",
			vmName: testVMName,
			vm:     testVM,
			memory: 1,
		},
		{
			name:      "ChangeMemoryNeedsForceToIncrease",
			args:      "--memory 10 --server test",
			vmName:    testVMName,
			vm:        testVM,
			memory:    10,
			shouldErr: true,
		},
		{
			name:   "ChangeCoresDoesNotNeedForceToDecrease",
			args:   "--cores 1 --server test",
			vmName: testVMName,
			vm:     testVM,
			cores:  1,
		},
		{
			name:      "ChangeCoresDoesNeedsForceToIncrease",
			args:      "--cores 4 --server test",
			vmName:    testVMName,
			vm:        testVM,
			cores:     4,
			shouldErr: true,
		},
		{
			name:   "EjectCdrom",
			args:   "--remove-cd --server test",
			vmName: testVMName,
			vm:     testVM,
			eject:  true,
		},
		{
			name:   "UpdateCdrom",
			args:   "--cd-url https://microsoft.com/windows.iso --server test",
			vmName: testVMName,
			vm:     testVM,
			cdrom:  "https://microsoft.com/windows.iso",
		},
		{
			name:   "EjectAndUpdateCdrom",
			args:   "--remove-cd --cd-url https://microsoft.com/windows.iso --server test",
			vmName: testVMName,
			vm:     testVM,
			cdrom:  "https://microsoft.com/windows.iso",
			eject:  true,
		},
		{
			name:      "CombinedChanges",
			args:      "--new-name new --memory 1 --hw-profile foo --cores 1 --remove-cd --cd-url https://microsoft.com/windows.iso --server test",
			vmName:    testVMName,
			vm:        testVM,
			hwProfile: "foo",
			memory:    1,
			cores:     1,
			cdrom:  "https://microsoft.com/windows.iso",
			eject:  true,
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
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			config.When("GetVirtualMachine").Return(defVM)
			client.When("GetVirtualMachine", test.vmName).Return(test.vm).Times(1)
			config.When("Force").Return(true)

			if test.hwProfile != "" {
				client.When("SetVirtualMachineHardwareProfile", test.vmName, test.hwProfile, []bool{test.hwProfileLock}).Return(nil).Times(1)
			}

			if test.move.expected {
				client.When("MoveVirtualMachine", test.vmName, test.move.newName).Return(nil).Times(1)
			}

			if test.memory > 0 && !test.shouldErr {
				client.When("SetVirtualMachineMemory", test.vmName, test.memory*1024).Return(nil).Times(1)
			}

			if test.cores > 0 && !test.shouldErr {
				client.When("SetVirtualMachineCores", test.vmName, test.cores).Return(nil).Times(1)
			}

			if test.eject || test.cdrom != "" {
				client.When("SetVirtualMachineCDROM", test.vmName, test.cdrom).Return(nil).Times(1)
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
