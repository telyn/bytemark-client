package main

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"testing"
)

func TestScheduleBackups(t *testing.T) {
	type ScheduleTest struct {
		Args []string

		Name      lib.VirtualMachineName
		DiscLabel string
		Start     string
		Interval  int

		ShouldErr  bool
		ShouldCall bool
		CreateErr  error
		BaseTestFn func(*testing.T, bool) (*mocks.Config, *mocks.Client)
	}

	tests := []ScheduleTest{
		{
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: baseTestSetup,
		},
		{
			Args:       []string{"vm-name"},
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: baseTestSetup,
		},
		{
			Args:       []string{"vm-name", "disc-label"},
			Name:       lib.VirtualMachineName{"vm-name", "default", "default-account"},
			DiscLabel:  "disc-label",
			Start:      "00:00",
			Interval:   86400,
			ShouldCall: true,
			ShouldErr:  false,
			BaseTestFn: baseTestAuthSetup,
		},
		{
			ShouldCall: true,
			Args:       []string{"vm-name.group.account", "disc-label", "3600"},
			Name:       lib.VirtualMachineName{"vm-name", "group", "account"},
			DiscLabel:  "disc-label",
			Start:      "00:00",
			Interval:   3600,
			BaseTestFn: baseTestAuthSetup,
		},
		{
			Args:       []string{"--start", "thursday", "vm-name", "disc-label", "3235"},
			Name:       lib.VirtualMachineName{"vm-name", "default", "default-account"},
			DiscLabel:  "disc-label",
			Start:      "thursday",
			Interval:   3235,
			ShouldCall: true,
			ShouldErr:  true,
			CreateErr:  fmt.Errorf("intermittent failure"),
			BaseTestFn: baseTestAuthSetup,
		},
	}

	var i int
	var test ScheduleTest
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestScheduleBackups #%d panicked.\r\n%v", i, r)
		}
	}()

	for i, test = range tests {
		fmt.Println(i) // fmt.Println still works even when the test panics - unlike t.Log

		config, client := test.BaseTestFn(t, false)
		config.When("GetVirtualMachine").Return(defVM)

		retSched := brain.BackupSchedule{
			StartDate: test.Start,
			Interval:  test.Interval,
			ID:        3442,
		}
		if test.ShouldCall {
			client.When("CreateBackupSchedule", test.Name, test.DiscLabel, test.Start, test.Interval).Return(retSched, test.CreateErr).Times(1)
		} else {
			client.When("CreateBackupSchedule", test.Name, test.DiscLabel, test.Start, test.Interval).Return(retSched, test.CreateErr).Times(0)
		}
		err := global.App.Run(append([]string{"bytemark", "schedule", "backups"}, test.Args...))
		checkErr(t, "TestScheduleBackups", i, test.ShouldErr, err)
		verifyAndReset(t, "TestScheduleBackups", i, client)
	}
}
