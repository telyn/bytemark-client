package main

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestUnscheduleBackups(t *testing.T) {

	tests := []struct {
		Args []string

		Name      lib.VirtualMachineName
		DiscLabel string
		ID        int

		ShouldErr  bool
		ShouldCall bool
		CreateErr  error
		BaseTestFn func(*testing.T, bool) (*mocks.Config, *mocks.Client, *cli.App)
	}{
		{
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: baseTestSetup,
		},
		{
			Args:       []string{"vm-name"},
			Name:       lib.VirtualMachineName{"vm-name", "default", "default-account"},
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: baseTestSetup,
		},
		{
			Args:       []string{"vm-name", "disc-label"},
			Name:       lib.VirtualMachineName{"vm-name", "default", "default-account"},
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: baseTestSetup,
		},
		{
			ShouldCall: true,
			Args:       []string{"vm-name", "disc-label", "324"},
			Name:       lib.VirtualMachineName{"vm-name", "default", "default-account"},
			DiscLabel:  "disc-label",
			ID:         324,
			BaseTestFn: baseTestAuthSetup,
		},
	}

	for i, test := range tests {
		config, client, app := test.BaseTestFn(t, false)
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
