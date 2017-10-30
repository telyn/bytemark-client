package admin_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
)

func TestUpdateBillingDefinition(t *testing.T) {
	tests := []struct {
		Command   string
		Expected  interface{}
		ShouldErr bool
	}{
		{
			Command:   "bytemark update billing-definition",
			ShouldErr: true,
		}, {
			Command: "bytemark update billing-definition --name trial_pence --value 2000",
			Expected: billing.Definition{
				Name:  "trial_pence",
				Value: "2000",
			},
		}, {
			Command: "bytemark update billing-definition --name trial_pence --value 2000 --group senior",
			Expected: billing.Definition{
				Name:           "trial_pence",
				Value:          "2000",
				UpdateGroupReq: "senior",
			},
		},
	}
	for i, test := range tests {
		_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)
		c.MockRequest = &mocks.Request{
			T:          t,
			StatusCode: 200,
		}

		err := app.Run(strings.Split(test.Command, " "))
		if err != nil && !test.ShouldErr {
			t.Errorf("TestUpdateBillingDefinition %d ERR: %s", i, err)
		} else if err == nil && test.ShouldErr {
			t.Errorf("TestUpdateBillingDefinition %d didn't err but should've", i)
		}
		assert.Equal(t, testutil.Name(i), test.Expected, c.MockRequest.RequestObject)
		if ok, err := c.Verify(); !ok && !test.ShouldErr {
			t.Fatal(err)
		}
	}

}

func TestUpdateVMMigrationWithSpeedAndDowntime(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm123", Group: "group", Account: "account"}
	speed := int64(8500000000000)
	downtime := 15
	c.When("UpdateVMMigration", vmName, &speed, &downtime).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "vm", "migration", "vm123.group.account", "8500000000000", "15"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateVMMigrationError(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	err := app.Run([]string{"bytemark", "update", "vm", "migration", "vm124.group.account"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateHead(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	options := lib.UpdateHead{}

	c.When("UpdateHead", "1", options).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "head", "1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateHeadError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	options := lib.UpdateHead{}

	c.When("UpdateHead", "1", options).Return(fmt.Errorf("Could not update head")).Times(1)

	err := app.Run([]string{"bytemark", "update", "head", "1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateTail(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	options := lib.UpdateTail{}

	c.When("UpdateTail", "1", options).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "tail", "1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateTailError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	options := lib.UpdateTail{}

	c.When("UpdateTail", "1", options).Return(fmt.Errorf("Could not update tail")).Times(1)

	err := app.Run([]string{"bytemark", "update", "tail", "1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateStoragePool(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	options := lib.UpdateStoragePool{}

	c.When("UpdateStoragePool", "1", options).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "storage", "pool", "1"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateStoragePoolError(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	options := lib.UpdateStoragePool{}

	c.When("UpdateStoragePool", "1", options).Return(fmt.Errorf("Could not update storage pool")).Times(1)

	err := app.Run([]string{"bytemark", "update", "storage", "pool", "1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
