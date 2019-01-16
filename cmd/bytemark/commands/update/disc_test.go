package update_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
)

func TestUpdateDisc(t *testing.T) {
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
		{
			name:      "RelativeSize",
			args:      "--force --server test --disc vda --new-size +10",
			vmName:    testVMName,
			discLabel: "vda",
			disc:      testDisc,
			newSize:   20,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
			config.When("Force").Return(true)
			config.When("GetVirtualMachine").Return(defVM)
			client.When("GetDisc", test.vmName, test.discLabel).Return(test.disc).Times(1)

			if !test.shouldErr {
				client.When("ResizeDisc", test.vmName, test.discLabel, test.newSize*1024).Return(nil).Times(1)
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

func TestUpdateDiscServer(t *testing.T) {
	defVM := lib.VirtualMachineName{
		Group:   "default",
		Account: "default-account",
	}
	oldVMName := lib.VirtualMachineName{
		VirtualMachine: "old-vm",
		Group:          "default",
		Account:        "default-account",
	}
	newVMName := lib.VirtualMachineName{
		VirtualMachine: "new-vm",
		Group:          "default",
		Account:        "default-account",
	}
	newVM := brain.VirtualMachine{
		ID: 999,
	}
	testDisc := brain.Disc{
		Label:        "vda",
		Size:         10240,
		StorageGrade: "sata",
	}

	args := strings.Split("bytemark update disc --force --server old-vm --disc vda --new-server new-vm", " ")

	is := is.New(t)

	t.Run("withServer", func(t *testing.T) {
		config, client, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)

		config.When("GetVirtualMachine").Return(defVM)
		client.When("GetVirtualMachine", newVMName).Return(newVM).Times(1)
		client.When("GetDisc", oldVMName, testDisc.Label).Return(testDisc).Times(1)

		url := "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s"
		urlValues := []string{"default-account", "default", "old-vm", "vda"}
		mockReturn := mocks.Request{
			T:              t,
			StatusCode:     200,
			ResponseObject: nil,
		}
		client.When("BuildRequest", "PUT", lib.BrainEndpoint, url, urlValues).Return(&mockReturn).Times(1)

		err := app.Run(args)
		is.Nil(err)
		if ok, err := client.Verify(); !ok {
			t.Fatal(err)
		}
	})
}
