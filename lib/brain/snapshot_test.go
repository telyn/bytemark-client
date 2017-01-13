package brain_test

import (
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"testing"
)

// TODO(telyn): all the prettyprint tests could be replaced with

func TestSnapshotPrettyPrint(t *testing.T) {
	tests := []prettyprint.Test{
		{
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
		{
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
		{
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
	}
	prettyprint.RunTests(t, "TestSnapshotPrettyPrint", tests)
}

func TestSnapshotsPrettyPrint(t *testing.T) {
	snapshots := brain.Snapshots{
		{
			Disc: brain.Disc{
				Label:        "kendrick-lamarr-to-pimp-a-butterfly",
				StorageGrade: "sata",
			},
		}, {
			Disc: brain.Disc{
				Label:        "kel-valhaal-new-introductory-lectures-on-transcendental-qabala",
				StorageGrade: "iceberg",
			},
		}, {
			Disc: brain.Disc{
				Label:        "dimmu-borgir-stormblåst",
				StorageGrade: "iceberg",
			},
		},
	}
	prettyprint.RunTests(t, "TestSnapshotsPrettyPrint", []prettyprint.Test{
		{
			Object: snapshots,
			Detail: prettyprint.Full,
			Expected: `kendrick-lamarr-to-pimp-a-butterfly (in progress)
kel-valhaal-new-introductory-lectures-on-transcendental-qabala
dimmu-borgir-stormblåst
`,
		}, {
			Object:   snapshots,
			Detail:   prettyprint.SingleLine,
			Expected: "3 snapshots",
		}, {
			Object: snapshots,
			Detail: prettyprint.Medium,
			Expected: `kendrick-lamarr-to-pimp-a-butterfly (in progress)
kel-valhaal-new-introductory-lectures-on-transcendental-qabala
dimmu-borgir-stormblåst
`,
		},
	})
}
