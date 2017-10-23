package lib

import (
	"bytes"
	"net"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/cheekybits/is"
)

func TestFormatAccount(t *testing.T) {
	is := is.New(t)
	b := new(bytes.Buffer)

	ip := net.IPv4(127, 0, 0, 1)
	gp := brain.Group{
		ID:   1,
		Name: "default",
		VirtualMachines: []brain.VirtualMachine{{
			Name:    "valid-vm",
			GroupID: 1,

			Autoreboot:            true,
			CdromURL:              "",
			Cores:                 1,
			Memory:                1,
			PowerOn:               true,
			HardwareProfile:       "fake-hardwareprofile",
			HardwareProfileLocked: false,
			ZoneName:              "default",
			ID:                    1,
			ManagementAddress:     ip,
			Deleted:               false,
			Hostname:              "valid-vm.default.account.fake-endpoint.example.com",
			Head:                  "fakehead",
		}},
	}
	acc := Account{
		BillingID: 2402,
		Name:      "test-account",
		Owner: billing.Person{
			Username: "test-user",
		},
		TechnicalContact: billing.Person{
			Username: "test-user",
		},
		Groups: []brain.Group{
			gp,
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
