package admin_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/util"
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
		c.MockRequest.AssertRequestObjectEqual(test.Expected)
		if ok, err := c.Verify(); !ok && !test.ShouldErr {
			t.Fatal(err)
		}
	}

}

func TestUpdateServerMigrationWithSpeedAndDowntime(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm123", Group: "group", Account: "account"}
	speed := int64(8500000000000)
	downtime := 15
	c.When("UpdateVMMigration", vmName, &speed, &downtime).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "update", "server-migration", "vm123.group.account", "8500000000000", "15"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateServerMigrationError(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	err := app.Run([]string{"bytemark", "update",  "server-migration", "vm124.group.account"})

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

	options := brain.StoragePool{}

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

	options := brain.StoragePool{}

	c.When("UpdateStoragePool", "1", options).Return(fmt.Errorf("Could not update storage pool")).Times(1)

	err := app.Run([]string{"bytemark", "update", "storage", "pool", "1"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUpdateMigration(t *testing.T) {
	// is := is.New(t)
	tests := []struct {
		name          string
		input         string
		modifications brain.MigrationJobModification
		shouldErr     bool
		cancelAll     bool
	}{
		{
			name:      "UpdateMigrationWithoutSpecifyingAnything",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "UpdateMigrationWithJustID",
			input:     "--id 1",
			shouldErr: true,
		},
		{
			name:      "UpdateMigrationCancelAllAndChangePriority",
			input:     "--id 1 --cancel-all --priority 5",
			shouldErr: true,
		},
		{
			name:      "UpdateMigrationCancelAllAndCancelAPool",
			input:     "--id 1 --cancel-all --cancel-pool t3-sata2",
			shouldErr: true,
		},
		{
			name:  "UpdateMigrationPriority",
			input: "--id 1 --priority 10",
			modifications: brain.MigrationJobModification{
				Cancel: brain.MigrationJobLocations{
					Discs: []util.NumberOrString{},
					Pools: []util.NumberOrString{},
					Tails: []util.NumberOrString{}},
				Options: brain.MigrationJobOptions{
					Priority: 10,
				},
			},
		},
		{
			name:  "UpdateMigrationPriorityAndCancelling",
			input: "--id 1 --priority 10 --cancel-pool t1-archive1 --cancel-disc disc.sata-1.8912 --cancel-tail tail2",
			modifications: brain.MigrationJobModification{
				Cancel: brain.MigrationJobLocations{
					Discs: []util.NumberOrString{"disc.sata-1.8912"},
					Pools: []util.NumberOrString{"t1-archive1"},
					Tails: []util.NumberOrString{"tail2"}},
				Options: brain.MigrationJobOptions{
					Priority: 10,
				},
			},
		},
		{
			name:  "UpdateMigrationCancellingMultiplesOfEach",
			input: "--id 1 --cancel-pool t1-archive1 --cancel-pool 679 --cancel-pool 2001:41c8:50:2::3a7 --cancel-disc disc.sata-1.8912 --cancel-disc 798 --cancel-disc 2001:41c8:50:2::3a6 --cancel-tail tail2 --cancel-tail 2001:41c8:50:2::3a8",
			modifications: brain.MigrationJobModification{
				Cancel: brain.MigrationJobLocations{
					Discs: []util.NumberOrString{"disc.sata-1.8912", "798", "2001:41c8:50:2::3a6"},
					Pools: []util.NumberOrString{"t1-archive1", "679", "2001:41c8:50:2::3a7"},
					Tails: []util.NumberOrString{"tail2", "2001:41c8:50:2::3a8"}},
			},
		},
		{
			name:      "UpdateMigrationCancelAll",
			input:     "--id 1 --cancel-all",
			cancelAll: true,
		},
		{
			name:  "UpdateMigrationCancelDisc",
			input: "--id 1 --cancel-disc disc.sata-1.8912",
			modifications: brain.MigrationJobModification{
				Cancel: brain.MigrationJobLocations{
					Discs: []util.NumberOrString{"disc.sata-1.8912"},
					Pools: []util.NumberOrString{},
					Tails: []util.NumberOrString{}},
			},
		},
		{
			name:  "UpdateMigrationCancelPool",
			input: "--id 1 --cancel-pool t1-archive1",
			modifications: brain.MigrationJobModification{
				Cancel: brain.MigrationJobLocations{
					Discs: []util.NumberOrString{},
					Pools: []util.NumberOrString{"t1-archive1"},
					Tails: []util.NumberOrString{}},
			},
		},
		{
			name:  "UpdateMigrationCancelTail",
			input: "--id 1 --cancel-tail tail2",
			modifications: brain.MigrationJobModification{
				Cancel: brain.MigrationJobLocations{
					Discs: []util.NumberOrString{},
					Pools: []util.NumberOrString{},
					Tails: []util.NumberOrString{"tail2"}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, client, app := testutil.BaseTestAuthSetup(t, false, admin.Commands)

			putReq := &mocks.Request{
				T:          t,
				StatusCode: 200,
			}

			client.When("BuildRequest", "PUT", lib.Endpoint(1),
				"/admin/migration_jobs/%s", []string{"1"}).Return(putReq).Times(1)

			args := fmt.Sprintf("bytemark --admin update migration %s", test.input)
			err := app.Run(strings.Split(args, " "))
			if !test.shouldErr && err != nil {
				t.Errorf("shouldn't err, but did: %T{%s}", err, err.Error())
			} else if test.shouldErr && err == nil {
				t.Errorf("should err, but didn't")
			}
			if !test.shouldErr {
				if ok, err := client.Verify(); !ok {
					t.Fatal(err)
				}
				if test.cancelAll {
					putReq.AssertRequestObjectEqual(map[string]interface{}{
						"cancel": map[string]interface{}{
							"all": true,
						},
					})
				} else {
					putReq.AssertRequestObjectEqual(test.modifications)
				}

			}

		})
	}
}
