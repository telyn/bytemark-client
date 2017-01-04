package lib

import (
	"bytes"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"testing"
)

func TestFormatOverview(t *testing.T) {
	is := is.New(t)
	b := new(bytes.Buffer)

	gp := getFixtureGroup()
	vm := getFixtureVM()
	megaGroup := &brain.Group{
		VirtualMachines: []*brain.VirtualMachine{
			&vm, &vm, &vm, &vm,
			&vm, &vm, &vm, &vm,
			&vm, &vm, &vm, &vm,
			&vm, &vm, &vm, &vm,
			&vm, &vm, &vm, &vm,
		},
	}
	accs := []*Account{
		&Account{
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
		},
		&Account{
			BillingID: 2403,
			Name:      "test-account-2",
			Owner: &billing.Person{
				Username: "test-user",
			},
			TechnicalContact: &billing.Person{
				Username: "test-user",
			},
			Groups: []*brain.Group{
				megaGroup,
			},
		},
		&Account{
			Name: "test-unowned-account",
			Groups: []*brain.Group{
				&gp,
			},
		},
		&Account{
			BillingID: 2406,
		},
	}

	err := FormatOverview(b, accs, accs[0], "test-user")
	if err != nil {
		t.Fatal(err)
	}

	expctd := `You are 'test-user'

Accounts you own:
  • 2402 - test-account (this is your default account)
  • 2403 - test-account-2

Other accounts you can access:
  • test-unowned-account
  • 2406 - [no bigv account]

Your default account (2402 - test-account)
  • default - 1 server
    ▸ valid-vm.default (powered on) in Default

`
	actual := b.String()

	is.Equal(expctd, actual)

	b.Reset()
	accs = []*Account{
		&Account{
			Name: "test-unowned-account",
			Groups: []*brain.Group{
				&gp,
			},
		},
		&Account{
			BillingID: 2406,
		},
	}

	err = FormatOverview(b, accs, nil, "test-user")
	if err != nil {
		t.Fatal(err)
	}

	expctd = `You are 'test-user'

Accounts you can access:
  • test-unowned-account
  • 2406 - [no bigv account]

It was not possible to determine your default account. Please set one using bytemark config set account.

`
	actual = b.String()
	IsEqualString(t, expctd, actual)
}
