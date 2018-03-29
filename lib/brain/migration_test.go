package brain

import (
	"bytes"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

func TestFormatMigration(t *testing.T) {
	tests := []struct {
		in     Migration
		detail prettyprint.DetailLevel
		exp    string
	}{
		{
			in: Migration{
				ID:             123,
				DiscID:         1,
				TailID:         2,
				Port:           3000,
				CreatedAt:      "2018-03-29T08:54:50.198Z",
				UpdatedAt:      "2018-03-29T08:54:50.198Z",
				MigrationJobID: 77,
			},
			detail: prettyprint.Full,
			exp: ` ▸ 123
     migration_job_id: 77
     tail_id: 2
     disc_id: 1
     port: 3000
     created_at: 2018-03-29T08:54:50.198Z
     updated_at: 2018-03-29T08:54:50.198Z
`,
		},
	}
	for _, test := range tests {
		b := new(bytes.Buffer)
		err := test.in.PrettyPrint(b, test.detail)
		if err != nil {
			t.Error(err)
		}
		if b.String() != test.exp {
			t.Errorf("unexpected output: %s", b.String())
		}
	}
}
