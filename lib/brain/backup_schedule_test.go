package brain_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
)

func TestBackupSchedulePrettyPrint(t *testing.T) {
	prettyprint.RunTests(t, []prettyprint.Test{
		{
			Object: brain.BackupSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.SingleLine,
			Expected: "#0: Daily starting from 2017-01-11 10:00:00",
		},
		{
			Object: brain.BackupSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.Medium,
			Expected: "#0: Daily starting from 2017-01-11 10:00:00",
		},
		{
			Object: brain.BackupSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.Full,
			Expected: "#0: Daily starting from 2017-01-11 10:00:00",
		},
	})

}

func TestBackupSchedulesPrettyPrint(t *testing.T) {
	schedules := brain.BackupSchedules{
		brain.BackupSchedule{
			ID:        24,
			StartDate: "2017-03-03 5:00:00",
			Interval:  35,
			Capacity:  1,
		},
		brain.BackupSchedule{
			ID:        4902,
			StartDate: "2017-03-03 5:00:00",
			Interval:  999,
		},
		brain.BackupSchedule{
			ID:        655,
			StartDate: "2017-03-03 5:00:00",
			Interval:  3600,
			Capacity:  500,
		},
		brain.BackupSchedule{
			ID:        234,
			StartDate: "2017-01-11 10:00:00",
			Interval:  86400,
		},
	}
	prettyprint.RunTests(t, []prettyprint.Test{
		{
			Object:   schedules,
			Expected: "Backups are taken every 35, 999, 3600 & 86400 seconds",
			Detail:   prettyprint.SingleLine,
		},
		{
			Object:   schedules,
			Expected: "Backups are taken every 35, 999, 3600 & 86400 seconds",
			Detail:   prettyprint.Medium,
		},
		{
			Object: schedules,
			Expected: `• #24: Every 35 seconds starting from 2017-03-03 5:00:00
• #4902: Every 999 seconds starting from 2017-03-03 5:00:00
• #655: Hourly starting from 2017-03-03 5:00:00
• #234: Daily starting from 2017-01-11 10:00:00
`,
			Detail: prettyprint.Full,
		},
	})
}
