package update_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

func TestUpdateDiscr(t *testing.T) {
	defVM := lib.VirtualMachineName{
		Group:   "default",
		Account: "default-account",
	}
	testVMName := lib.VirtualMachineName{
		VirtualMachine: "test",
		Group:          "default",
		Account:        "default-account",
	}
	testDisc := brain.Disc{
		Size:         10240,
		StorageGrade: "sata",
	}
	tests := []struct {
		name      string
		args      string
		vmName    lib.VirtualMachineName
		discLabel string
		disc      brain.Disc
		newSize   int
		shouldErr bool
	}{
		{
			name:      "MissingForceFlag",
			args:      "--server test --disc vda --new-size 100",
			vmName:    testVMName,
			discLabel: "vda",
			disc:      testDisc,
			shouldErr: true,
		},
		{
			name:      "AbsoluteSize",
			args:      "--force --server test --disc vda --new-size 100",
			vmName:    testVMName,
			discLabel: "vda",
			disc:      testDisc,
			newSize:   100,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			config.When("Force").Return(true)
			config.When("GetVirtualMachine").Return(defVM)
			client.When("GetDisc", test.vmName, test.discLabel).Return(&test.disc).Times(1)

			if !test.shouldErr {
				client.When("ResizeDisc", test.vmName, test.newSize*1024).Return(nil).Times(1)
			}

			args := strings.Split("bytemark update disc "+test.args, " ")
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
