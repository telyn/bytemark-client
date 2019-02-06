package main

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestUnscheduleBackups(t *testing.T) {

	tests := []struct {
		Args []string

		Name      pathers.VirtualMachineName
		DiscLabel string
		ID        int

		ShouldErr  bool
		ShouldCall bool
		CreateErr  error
		BaseTestFn func(*testing.T, bool, []cli.Command) (*mocks.Config, *mocks.Client, *cli.App)
	}{
		{
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: testutil.BaseTestSetup,
		},
		{
			Args:       []string{"vm-name"},
			Name:       pathers.VirtualMachineName{VirtualMachine: "vm-name", GroupName: pathers.GroupName{Group: "default", Account: "default-account"}},
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: testutil.BaseTestSetup,
		},
		{
			Args:       []string{"vm-name", "disc-label"},
			Name:       pathers.VirtualMachineName{VirtualMachine: "vm-name", GroupName: pathers.GroupName{Group: "default", Account: "default-account"}},
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: testutil.BaseTestSetup,
		},
		{
			ShouldCall: true,
			Args:       []string{"vm-name", "disc-label", "324"},
			Name:       pathers.VirtualMachineName{VirtualMachine: "vm-name", GroupName: pathers.GroupName{Group: "default", Account: "default-account"}},
			DiscLabel:  "disc-label",
			ID:         324,
			BaseTestFn: testutil.BaseTestAuthSetup,
		},
	}

	for i, test := range tests {
		config, client, app := test.BaseTestFn(t, false, commands)
		config.When("GetVirtualMachine").Return(defVM)
		fmt.Println(i) // fmt.Println still works even when the test panics - unlike t.Log

		if test.ShouldCall {
			client.When("DeleteBackupSchedule", test.Name, test.DiscLabel, test.ID).Return(test.CreateErr).Times(1)
		} else {
			client.When("DeleteBackupSchedule", test.Name, test.DiscLabel, test.ID).Return(test.CreateErr).Times(0)
		}
		err := app.Run(append([]string{"bytemark", "unschedule", "backups"}, test.Args...))
		checkErr(t, "TestUnscheduleBackups", i, test.ShouldErr, err)
		verifyAndReset(t, "TestUnscheduleBackups", i, client)
	}
}
