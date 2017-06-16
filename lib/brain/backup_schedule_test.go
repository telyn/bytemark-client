package brain_test

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"testing"
)

func TestBackupSchedulePrettyPrint(t *testing.T) {
	prettyprint.RunTests(t, []prettyprint.Test{
		{
			Object: brain.BackupSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.SingleLine,
			Expected: "Every 86400 seconds starting from 2017-01-11 10:00:00",
		},
		{
			Object: brain.BackupSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.Medium,
			Expected: "Every 86400 seconds starting from 2017-01-11 10:00:00",
		},
		{
			Object: brain.BackupSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.Full,
			Expected: "Every 86400 seconds starting from 2017-01-11 10:00:00",
		},
	})

}

func TestBackupSchedulesPrettyPrint(t *testing.T) {
	schedules := brain.BackupSchedules{
		&brain.BackupSchedule{
			ID:        24,
			StartDate: "2017-03-03 5:00:00",
			Interval:  35,
		},
		&brain.BackupSchedule{
			ID:        4902,
			StartDate: "2017-03-03 5:00:00",
			Interval:  999,
		},
		&brain.BackupSchedule{
			ID:        655,
			StartDate: "2017-03-03 5:00:00",
			Interval:  3306,
		},
		&brain.BackupSchedule{
			ID:        234,
			StartDate: "2017-01-11 10:00:00",
			Interval:  86400,
		},
	}
	prettyprint.RunTests(t, []prettyprint.Test{
		{
			Object:   schedules,
			Expected: "Backups are taken every 35, 999, 3306 & 86400 seconds",
			Detail:   prettyprint.SingleLine,
		},
		{
			Object:   schedules,
			Expected: "Backups are taken every 35, 999, 3306 & 86400 seconds",
			Detail:   prettyprint.Medium,
		},
		{
			Object: schedules,
			Expected: `• #24 - Every 35 seconds starting from 2017-03-03 5:00:00
• #4902 - Every 999 seconds starting from 2017-03-03 5:00:00
• #655 - Every 3306 seconds starting from 2017-03-03 5:00:00
• #234 - Every 86400 seconds starting from 2017-01-11 10:00:00
`,
			Detail: prettyprint.Full,
		},
	})
}
