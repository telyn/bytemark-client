package admin_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
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

func TestAdminShowMigratingVMsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	vms := []brain.VirtualMachine{getFixtureVM()}
	c.When("GetMigratingVMs").Return(&vms, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show migrating vms", " "))
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

func TestAdminShowStoppedEligibleVMsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	vms := []brain.VirtualMachine{getFixtureVM()}
	c.When("GetStoppedEligibleVMs").Return(&vms, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show stopped eligible vms", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowRecentVMsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	vms := []brain.VirtualMachine{getFixtureVM()}
	c.When("GetRecentVMs").Return(&vms, nil).Times(1)

	err := app.Run(strings.Split("bytemark --admin show recent vms", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
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
