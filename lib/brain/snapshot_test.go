package brain_test

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"testing"
)

// TODO(telyn): all the prettyprint tests could be replaced with

func TestSnapshotPrettyPrint(t *testing.T) {
	prettyprint.RunTests(t, "TestSnapshotPrettyPrint", prettyprint.Tests{
		prettyprint.Test{
			Object: brain.Snapshot{
				Disc: brain.Disc{
					Label:        "taylorswift-1989-this-sick-beat-trademark-violations-20170101",
					Size:         25600,
					StorageGrade: "sata",
				},
			},
			Detail:   prettyprint.SingleLine,
			Expected: "taylorswift-1989-this-sick-beat-trademark-violations-20170101 (in progress)",
		},
		prettyprint.Test{
			Object: brain.Snapshot{
				Disc: brain.Disc{
					Label:        "taylorswift-1989-this-sick-beat-trademark-violations-20170101",
					Size:         25600,
					StorageGrade: "iceberg",
				},
			},
			Detail:   prettyprint.Medium,
			Expected: "taylorswift-1989-this-sick-beat-trademark-violations-20170101",
		},
		prettyprint.Test{
			Object: brain.Snapshot{
				Disc: brain.Disc{
					Label:        "taylorswift-1989-this-sick-beat-trademark-violations-20170101",
					Size:         25600,
					StorageGrade: "iceberg",
				},
			},
			Detail:   prettyprint.Full,
			Expected: "taylorswift-1989-this-sick-beat-trademark-violations-20170101",
		},
	})
}
