package lib

import (
	"bytes"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/prettyprint"
	"github.com/cheekybits/is"
	"testing"
)

func TestFormatAccount(t *testing.T) {
	is := is.New(t)
	b := new(bytes.Buffer)

	gp := getFixtureGroup()
	acc := Account{
		BillingID: 2402,
		Name:      "test-account",
		Owner: &billing.Person{
			Username: "test-user",
		},
		TechnicalContact: &billing.Person{
			Username: "test-user",
		},
		Groups: []*brain.Group{
			&gp,
		},
	}

	err := acc.PrettyPrint(b, prettyprint.Full)
	if err != nil {
		t.Error(err)
	}

	is.Equal(`2402 - test-account
  • default - 1 server
    ▸ valid-vm.default (powered on) in Default
`, b.String())
}
