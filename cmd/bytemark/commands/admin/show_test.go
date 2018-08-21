package admin_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
)

func TestAdminShowVLANsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	vlans := []brain.VLAN{getFixtureVLAN()}
	c.When("GetVLANs").Return(&vlans, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show vlans", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowVLANCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	vlanNum := 1
	vlan := getFixtureVLAN()
	c.When("GetVLAN", vlanNum).Return(&vlan, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show vlan 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowIPRangesCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	ipRanges := []brain.IPRange{getFixtureIPRange()}
	c.When("GetIPRanges").Return(&ipRanges, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show ip ranges", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowIPRangeCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	ipRangeID := "1"
	ipRange := getFixtureIPRange()
	c.When("GetIPRange", ipRangeID).Return(&ipRange, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show ip range 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowIPRangeWithIPRangeCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	rangeString := "192.168.33.0/24"
	ipRange := getFixtureIPRange()
	c.When("GetIPRange", rangeString).Return(&ipRange, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show ip range 192.168.33.0/24", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowHeadsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	heads := []brain.Head{getFixtureHead()}
	c.When("GetHeads").Return(&heads, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show heads", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowHeadCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	headID := "1"
	head := getFixtureHead()
	c.When("GetHead", headID).Return(&head, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show head 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowTailsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	tails := []brain.Tail{getFixtureTail()}
	c.When("GetTails").Return(&tails, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show tails", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowTailCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	tailID := "1"
	tail := getFixtureTail()
	c.When("GetTail", tailID).Return(&tail, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show tail 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowStoragePoolsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	storagePool := []brain.StoragePool{getFixtureStoragePool()}
	c.When("GetStoragePools").Return(&storagePool, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show storage pools", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowStoragePoolCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	storagePoolID := "1"
	storagePool := getFixtureStoragePool()
	c.When("GetStoragePool", storagePoolID).Return(&storagePool, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show storage pool 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowMigratingDiscsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	discs := brain.Discs{{ID: 134, StorageGrade: "sata", Size: 25600}}
	c.When("GetMigratingDiscs").Return(discs, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show migrating discs", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowMigrations(t *testing.T) {
	tests := []struct {
		name string
		mjs  brain.MigrationJobs
	}{
		{
			name: "NoJobs",
			mjs:  brain.MigrationJobs{},
		},
		{
			name: "OneJob",
			mjs: brain.MigrationJobs{
				{ID: 123},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			_, client, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)
			client.When("BuildRequest", "GET", lib.BrainEndpoint, "/admin/migration_jobs?unfinished=1%s", []string{""}).Return(&mocks.Request{
				T:              t,
				StatusCode:     200,
				ResponseObject: test.mjs,
			}).Times(1)
			err := app.Run(strings.Split("bytemark --admin show migrations", " "))
			is.Nil(err)
			if ok, err := client.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}

func TestAdminShowMigration(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		shouldErr bool
	}{
		{
			name:      "MissingID",
			shouldErr: true,
		},
		{
			name: "Successful",
			id:   123,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, client, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)
			args := []string{"bytemark", "--admin", "show", "migration"}
			if test.id > 0 {
				client.When("BuildRequest", "GET", lib.BrainEndpoint, "/admin/migration_jobs/%s", []string{strconv.Itoa(test.id)}).Return(&mocks.Request{
					T:              t,
					StatusCode:     200,
					ResponseObject: brain.MigrationJob{ID: test.id},
				}).Times(1)
				client.When("BuildRequest", "GET", lib.BrainEndpoint, "/admin/migration_jobs/%s/migrations", []string{strconv.Itoa(test.id)}).Return(&mocks.Request{
					T:              t,
					StatusCode:     200,
					ResponseObject: brain.Migrations{},
				}).Times(1)
				args = append(args, "--id", strconv.Itoa(test.id))
			}
			err := app.Run(args)
			if !test.shouldErr && err != nil {
				t.Errorf("shouldn't err, but did: %T{%s}", err, err.Error())
			} else if test.shouldErr && err == nil {
				t.Errorf("should err, but didn't")
			}
			if !test.shouldErr {
				if ok, err := client.Verify(); !ok {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestAdminShowDiscByIDCommand(t *testing.T) {
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	discID := 132
	disc := getFixtureDisc()
	c.When("GetDiscByID", discID).Return(&disc, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show disc by id 132", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowDependantServers(t *testing.T) {
	tests := []struct {
		name	string
		url		string
		args 	[]string
		values 	[]string
	}{
		{
			name: "head",
			url: "/admin/heads/%s/virtual_machines",
			args: []string{"--head", "123"},
			values: []string{"123"},
		},
		{
			name: "tail",
			url: "/admin/tails/%s/virtual_machines",
			args: []string{"--tail", "123"},
			values: []string{"123"},
		},
		{
			name: "storage pool",
			url: "/admin/storage_pools/%s/virtual_machines",
			args: []string{"--storage-pool", "123"},
			values: []string{"123"},
		},
		{
			name: "head",
			url: "/admin/heads/%s/virtual_machines?at=%s",
			args: []string{"--head", "123", "--at", "2018-08-21T15:00:00+0000"},
			values: []string{"123", "2018-08-21T15:00:00+0000"},
		},
		{
			name: "tail",
			url: "/admin/tails/%s/virtual_machines?at=%s",
			args: []string{"--tail", "123", "--at", "2018-08-21T15:00:00+0000"},
			values: []string{"123", "2018-08-21T15:00:00+0000"},
		},
		{
			name: "storage pool",
			url: "/admin/storage_pools/%s/virtual_machines?at=%s",
			args: []string{"--storage-pool", "123", "--at", "2018-08-21T15:00:00+0000"},
			values: []string{"123", "2018-08-21T15:00:00+0000"},
		},
	}

	baseArgs := []string{"bytemark", "--admin", "show", "dependant", "servers"}
	servers := mocks.Request{
		T:              t,
		StatusCode:     200,
		ResponseObject: []brain.VirtualMachine{
			{ID: 1, Name: "Test1"},
			{ID: 2, Name: "Test2"},
		},
	}

	is := is.New(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, client, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

			client.When("BuildRequest", "GET", lib.BrainEndpoint, test.url, test.values).Return(&servers).Times(1)

			err := app.Run(append(baseArgs, test.args...))
			is.Nil(err)
			if ok, err := client.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}

func TestAdminShowDependantDiscs(t *testing.T) {
	tests := []struct {
		name	string
		url		string
		args 	[]string
		values 	[]string
	}{
		{
			name: "tail",
			url: "/admin/tails/%s/discs",
			args: []string{"--tail", "123"},
			values: []string{"123"},
		},
		{
			name: "storage pool",
			url: "/admin/storage_pools/%s/discs",
			args: []string{"--storage-pool", "123"},
			values: []string{"123"},
		},
		{
			name: "tail",
			url: "/admin/tails/%s/discs?at=%s",
			args: []string{"--tail", "123", "--at", "2018-08-21T15:00:00+0000"},
			values: []string{"123", "2018-08-21T15:00:00+0000"},
		},
		{
			name: "storage pool",
			url: "/admin/storage_pools/%s/discs?at=%s",
			args: []string{"--storage-pool", "123", "--at", "2018-08-21T15:00:00+0000"},
			values: []string{"123", "2018-08-21T15:00:00+0000"},
		},
	}

	baseArgs := []string{"bytemark", "--admin", "show", "dependant", "discs"}
	discs := mocks.Request{
		T:              t,
		StatusCode:     200,
		ResponseObject: []brain.Disc{
			{ID: 1},
			{ID: 2},
		},
	}

	is := is.New(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, client, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

			client.When("BuildRequest", "GET", lib.BrainEndpoint, test.url, test.values).Return(&discs).Times(1)

			err := app.Run(append(baseArgs, test.args...))
			is.Nil(err)
			if ok, err := client.Verify(); !ok {
				t.Fatal(err)
			}
		})
	}
}
