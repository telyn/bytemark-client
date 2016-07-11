package lib

import (
	"bytes"
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
		Owner: &Person{
			Username: "test-user",
		},
		TechnicalContact: &Person{
			Username: "test-user",
		},
		Groups: []*Group{
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
	megaGroup := &Group{
		VirtualMachines: []*VirtualMachine{
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
			Owner: &Person{
				Username: "test-user",
			},
			TechnicalContact: &Person{
				Username: "test-user",
			},
			Groups: []*Group{
				&gp,
			},
		},
		&Account{
			BillingID: 2403,
			Name:      "test-account-2",
			Owner: &Person{
				Username: "test-user",
			},
			TechnicalContact: &Person{
				Username: "test-user",
			},
			Groups: []*Group{
				megaGroup,
			},
		},
		&Account{
			Name: "test-unowned-account",
			Groups: []*Group{
				&gp,
			},
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

Your default account (2402 - test-account)
  • default - 1 server
    ▸ valid-vm.default (powered on) in Default

`
	actual := b.String()

	//	expctdjs, err := json.Marshal(expctd)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	actualjs, err := json.Marshal(actual)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	expr := []rune(expctd)
	//	actr := []rune(actual)
	//
	//	is.Equal(len(expr), len(actr))
	//
	//	for i := 0; i < len(expr); i++ {
	//		if expr[i] != actr[i] {
	//			fmt.Printf("chr #%d differs. e:'%c' a:'%c'\r\n", i, expr[i], actr[i])
	//		}
	//	}

	//	fmt.Printf("\r\n%s\r\n%s", map[string]string{"data": string(expctdjs)}, map[string]string{"data": string(actualjs)})
	is.Equal(expctd, actual)
}
