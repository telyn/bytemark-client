package brain_test

import (
	"bytes"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

func TestAPIKeyPrettyPrint(t *testing.T) {
	tests := []struct {
		name  string
		level prettyprint.DetailLevel
		in    brain.APIKey
		out   string
	}{
		{
			name:  "single line",
			level: prettyprint.SingleLine,
			in: brain.APIKey{
				Label: "jeff",
			},
			out: "jeff",
		}, {
			name:  "single line (expired)",
			level: prettyprint.SingleLine,
			in: brain.APIKey{
				Label:     "jeff",
				ExpiresAt: "2006-01-01T01:01:01.000-0000",
			},
			out: "jeff (expired)",
		}, {
			name:  "full no privs expired",
			level: prettyprint.Full,
			in: brain.APIKey{
				Label:     "jeff",
				ExpiresAt: "2006-01-01T01:01:01.0124-0000",
			},
			out: `jeff (expired)
  Expired: 2006-01-01T01:01:01.0124-0000
`,
		}, {
			name:  "full with privs",
			level: prettyprint.Full,
			in: brain.APIKey{
				Label:     "jeff",
				ExpiresAt: "3000-01-01T01:01:01-0000",
				Privileges: brain.Privileges{
					{
						Username:    "jeffathan",
						AccountID:   23,
						AccountName: "jeffadiah",
						Level:       "account_admin",
						APIKeyID:    4,
					},
				},
			},
			out: `jeff
  Expires: 3000-01-01T01:01:01-0000

  Privileges:
    * account_admin on account jeffadiah for jeffathan
`,
		}, {
			name:  "full with key",
			level: prettyprint.Full,
			in: brain.APIKey{
				Label:     "jeff",
				APIKey:    "abcdefgh",
				ExpiresAt: "3006-01-01T01:01:01-0000",
			},
			out: `jeff
  Expires: 3006-01-01T01:01:01-0000
  Key: apikey.abcdefgh
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			err := test.in.PrettyPrint(&buf, test.level)
			if err != nil {
				t.Fatal(err)
			}
			str := buf.String()
			if str != test.out {
				t.Errorf("Output didn't match expected\nexpected: %q\n  actual: %q", test.out, str)
			}
		})
	}

}
