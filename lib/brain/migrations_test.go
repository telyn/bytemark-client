package brain

import (
	"bytes"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

func TestFormatMigrations(t *testing.T) {
	tests := []struct {
		name   string
		in     Migrations
		detail prettyprint.DetailLevel
		exp    string
	}{
		{
			name:   "NoMigrations",
			in:     Migrations{},
			detail: prettyprint.Full,
			exp:    ``,
		},
		{
			name: "OneMigration",
			in: Migrations{{
				ID:             123,
				DiscID:         1,
				TailID:         2,
				Port:           3000,
				CreatedAt:      "2018-03-29T08:54:50.198Z",
				UpdatedAt:      "2018-03-29T08:54:50.198Z",
				MigrationJobID: 77,
			}},
			detail: prettyprint.Full,
			exp:    `     ▸ 123 disc: 1` + "\n",
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
