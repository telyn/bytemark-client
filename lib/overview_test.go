package lib

import (
	"bytes"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"testing"
)

func TestFormatOverview(t *testing.T) {
	b := new(bytes.Buffer)

	gp := getFixtureGroup()
	vm := getFixtureVM()
	megaGroup := &brain.Group{
		Name: "mega-group",
		VirtualMachines: []*brain.VirtualMachine{
			&vm, &vm, &vm, &vm,
			&vm, &vm, &vm, &vm,
			&vm, &vm, &vm, &vm,
			&vm, &vm, &vm, &vm,
			&vm, &vm, &vm, &vm,
		},
	}
	tests := []struct {
		Accounts       []*Account
		DefaultAccount *Account
		Expected       string
	}{
		{
			Accounts: []*Account{
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
					IsDefaultAccount: true,
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
			},
			Expected: `You are 'test-user'

Accounts you own:
  • 2402 - test-account (this is your default account)
  • 2403 - test-account-2

Other accounts you can access:
  • test-unowned-account
  • 2406 - [no bigv account]

Your default account (2402 - test-account)
  • default - 1 server
    ▸ valid-vm.default (powered on) in Default

`,
		}, {
			Accounts: []*Account{
				&Account{
					Name: "test-unowned-account",
					Groups: []*brain.Group{
						&gp,
					},
				},
				&Account{
					BillingID: 2406,
				},
			},
			Expected: `You are 'test-user'

Accounts you can access:
  • test-unowned-account
  • 2406 - [no bigv account]

It was not possible to determine your default account. Please set one using bytemark config set account.

`,
		}, {
			Accounts: []*Account{
				&Account{
					Name: "test-account",
					Groups: []*brain.Group{
						&brain.Group{
							Name: "default",
							VirtualMachines: []*brain.VirtualMachine{
								&vm, &vm, &vm, &vm, &vm,
							},
						},
						megaGroup,
					},
					IsDefaultAccount: true,
				},
			},
			Expected: `You are 'test-user'

Accounts you can access:
  • test-account (this is your default account)

Your default account (test-account)
  • default - 5 servers
    ▸ valid-vm.default (powered on) in Default
    ▸ valid-vm.default (powered on) in Default
    ▸ valid-vm.default (powered on) in Default
    ▸ valid-vm.default (powered on) in Default
    ▸ valid-vm.default (powered on) in Default
  • mega-group - 20 servers

`,
		},
	}

	for i, test := range tests {
		err := FormatOverview(b, test.Accounts, nil, "test-user")
		if err != nil {
			t.Fatal(err)
		}

		actual := b.String()
		if test.Expected != actual {
			t.Errorf("TestFormatOverview %d FAIL\r\nexpected %s\r\nreceived %s", i, test.Expected, actual)
		}

		b.Reset()
	}
}
