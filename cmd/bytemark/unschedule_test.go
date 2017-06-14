package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"testing"
)

func TestUnscheduleBackups(t *testing.T) {
	config, client := baseTestSetup(t, false)
	config.When("Get", "token").Return("test-token")
	config.When("GetIgnoreErr", "yubikey").Return("")
	config.When("GetVirtualMachine").Return(&defVM)

	tests := []struct {
		Args []string

		Name      lib.VirtualMachineName
		DiscLabel string
		ID        int

		ShouldErr  bool
		ShouldCall bool
		CreateErr  error
	}{
		{
			ShouldCall: false,
			ShouldErr:  true,
		},
		{
			Args:       []string{"vm-name"},
			Name:       lib.VirtualMachineName{"vm-name", "default", "default-account"},
			ShouldCall: false,
			ShouldErr:  true,
		},
		{
			Args:       []string{"vm-name", "disc-label"},
			Name:       lib.VirtualMachineName{"vm-name", "default", "default-account"},
			ShouldCall: false,
			ShouldErr:  true,
		},
		{
			ShouldCall: true,
			Args:       []string{"vm-name", "disc-label", "324"},
			Name:       lib.VirtualMachineName{"vm-name", "default", "default-account"},
			DiscLabel:  "disc-label",
			ID:         324,
		},
	}

	for i, test := range tests {
		fmt.Println(i) // fmt.Println still works even when the test panics - unlike t.Log
		client.When("AuthWithToken", "test-token").Return(nil)

		if test.ShouldCall {
			client.When("DeleteBackupSchedule", test.Name, test.DiscLabel, test.ID).Return(test.CreateErr).Times(1)
		} else {
			client.When("DeleteBackupSchedule", test.Name, test.DiscLabel, test.ID).Return(test.CreateErr).Times(0)
		}
		err := global.App.Run(append([]string{"bytemark", "unschedule", "backups"}, test.Args...))
		checkErr(t, "TestUnscheduleBackups", i, test.ShouldErr, err)
		verifyAndReset(t, "TestUnscheduleBackups", i, client)
	}
}
