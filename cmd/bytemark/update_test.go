package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
)

func TestUpdateBmbilling(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		Command  string
		Expected interface{}
	}{
		{
			Command:  "bytemark update bmbilling --trial-days 7",
			Expected: billing.Definitions{TrialDays: 7},
		}, {
			Command:  "bytemark update bmbilling --trial-pence 2000",
			Expected: billing.Definitions{TrialPence: 2000},
		}, {
			Command: "bytemark update bmbilling --trial-days 7 --trial-pence 2000",
			Expected: billing.Definitions{
				TrialDays:  7,
				TrialPence: 2000,
			},
		},
	}
	for _, test := range tests {
		_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)
		c.MockRequest = &mocks.Request{
			T:          t,
			StatusCode: 200,
		}

		err := app.Run(strings.Split(test.Command, " "))
		is.Nil(err)
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}

}

func TestUpdateVMMigrationWithSpeedAndDowntime(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

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
	config, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

	config.When("GetVirtualMachine").Return(defVM)

	err := app.Run([]string{"bytemark", "update", "vm", "migration", "vm124.group.account"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateHead(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

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
	_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

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
	_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

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
	_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

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
	_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

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
	_, c, app := testutil.BaseTestAuthSetup(t, true, adminCommands)

	options := lib.UpdateStoragePool{}

	c.When("UpdateStoragePool", "1", options).Return(fmt.Errorf("Could not update storage pool")).Times(1)

	err := app.Run([]string{"bytemark", "update", "storage", "pool", "1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
