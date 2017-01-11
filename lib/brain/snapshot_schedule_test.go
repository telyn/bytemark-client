package brain_test

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"testing"
)

func TestSnapshotSchedulePrettyPrint(t *testing.T) {
	prettyprint.RunTests(t, "TestSnapshotSchedulePrettyPrint", []prettyprint.Test{
		{
			Object: brain.SnapshotSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.SingleLine,
			Expected: "Every 86400 seconds starting from 2017-01-11 10:00:00",
		},
		{
			Object: brain.SnapshotSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.Medium,
			Expected: "Every 86400 seconds starting from 2017-01-11 10:00:00",
		},
		{
			Object: brain.SnapshotSchedule{
				StartDate: "2017-01-11 10:00:00",
				Interval:  86400,
			},
			Detail:   prettyprint.Full,
			Expected: "Every 86400 seconds starting from 2017-01-11 10:00:00",
		},
	})

}
