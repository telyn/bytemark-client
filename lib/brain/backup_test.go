package brain_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// TODO(telyn): all the prettyprint tests could be replaced with

func TestBackupPrettyPrint(t *testing.T) {
	tests := []prettyprint.Test{
		{
			Object: brain.Backup{
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
			Object: brain.Backup{
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
			Object: brain.Backup{
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
	prettyprint.RunTests(t, tests)
}

func TestBackupsPrettyPrint(t *testing.T) {
	backups := brain.Backups{
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
	prettyprint.RunTests(t, []prettyprint.Test{
		{
			Object: backups,
			Detail: prettyprint.Full,
			Expected: `kendrick-lamarr-to-pimp-a-butterfly (in progress)
kel-valhaal-new-introductory-lectures-on-transcendental-qabala
dimmu-borgir-stormblåst
`,
		}, {
			Object:   backups,
			Detail:   prettyprint.SingleLine,
			Expected: "3 backups",
		}, {
			Object: backups,
			Detail: prettyprint.Medium,
			Expected: `kendrick-lamarr-to-pimp-a-butterfly (in progress)
kel-valhaal-new-introductory-lectures-on-transcendental-qabala
dimmu-borgir-stormblåst
`,
		},
	})
}
