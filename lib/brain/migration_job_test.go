package brain

import (
	"bytes"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

func TestFormatMigrationJob(t *testing.T) {
	tests := []struct {
		name   string
		in     MigrationJob
		detail prettyprint.DetailLevel
		exp    string
	}{
		{
			name: "SingleLine",
			in: MigrationJob{
				ID: 456,
				Queue: MigrationJobQueue{
					Discs: []int{100, 101},
				},
			},
			detail: prettyprint.SingleLine,
			exp:    ` ▸ 456 queue: 100 101`,
		},
		{
			name: "FullDetail",
			in: MigrationJob{
				ID: 123,
				Queue: MigrationJobQueue{
					Discs: []int{1, 2},
				},
				Status: MigrationJobStatus{
					Discs: MigrationJobDiscStatus{
						Done:    []int{3},
						Skipped: []int{5},
					},
				},
				Priority: 10,
			},
			detail: prettyprint.Full,
			exp: ` ▸ 123
   queue: 
     discs:
       • 1
       • 2
   status: 
     done:
       • 3
     skipped:
       • 5
   priority: 10
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			err := test.in.PrettyPrint(b, test.detail)
			if err != nil {
				t.Error(err)
			}
			if b.String() != test.exp {
				t.Errorf("unexpected output: %s", b.String())
			}
		})
	}
}
