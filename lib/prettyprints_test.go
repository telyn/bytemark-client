package lib

import (
	"bytes"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"testing"
)

func TestFormatVM(t *testing.T) {
	is := is.New(t)
	b := new(bytes.Buffer)
	vm, _, _ := getFixtureVMWithManyIPs()

	err := FormatVirtualMachine(b, &vm, "server_summary")
	if err != nil {
		t.Error(err)
	}

	is.Equal(" ▸ valid-vm.default (powered on) in Default", b.String())
	b.Truncate(0)
	err = FormatVirtualMachine(b, &vm, "server_spec")
	if err != nil {
		t.Error(err)
	}
	is.Equal("   192.168.1.16 - 1 core, 1MiB, 25GiB on 1 disc", b.String())
}

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

	err := FormatAccount(b, &acc, &Account{Name: ""}, "account_name")
	if err != nil {
		t.Error(err)
	}
	is.Equal(`2402 - test-account`, b.String())

	b.Truncate(0)

	err = FormatAccount(b, &acc, &Account{Name: ""}, "account_overview")
	if err != nil {
		t.Error(err)
	}

	is.Equal(`2402 - test-account
  • default - 1 server
    ▸ valid-vm.default (powered on) in Default
`, b.String())
}

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
