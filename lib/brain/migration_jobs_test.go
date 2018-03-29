package brain

import (
	"bytes"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

func TestFormatMigrationJobs(t *testing.T) {
	tests := []struct {
		name   string
		in     MigrationJobs
		detail prettyprint.DetailLevel
		exp    string
	}{
		{
			name:   "NoJobs",
			in:     MigrationJobs{},
			detail: prettyprint.Full,
			exp:    ``,
		},
		{
			name: "OneJob",
			in: MigrationJobs{{
				ID: 123,
				Queue: MigrationJobQueue{
					Discs: []int{1, 2},
				},
			}},
			detail: prettyprint.Full,
			exp:    ` ▸ 123 queue: 1 2` + "\n",
		},
		{
			name: "TwoJobs",
			in: MigrationJobs{
				{
					ID: 123,
					Queue: MigrationJobQueue{
						Discs: []int{1, 2},
					},
				},
				{
					ID: 456,
					Queue: MigrationJobQueue{
						Discs: []int{3, 4},
					},
				},
			},
			detail: prettyprint.Full,
			exp:    ` ▸ 123 queue: 1 2` + "\n" + ` ▸ 456 queue: 3 4` + "\n",
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
