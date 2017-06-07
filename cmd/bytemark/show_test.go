package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
	"runtime/debug"
	"strings"
	"testing"
)

type showAccountTest struct {
	// the command line to call
	Input string
	// the account name that's in the config, either because it's what was set as a global flag or it's in the config dir
	ConfigAccount string
	// the account that should be attempted to be got - if blank, expect GetDefaultAccount
	AccountToGet string
	ShouldErr    bool
}

func TestShowAccountCommand(t *testing.T) {
	baseShowAccountSetup := func(c *mocks.Client, config *mocks.Config, configAccount string) {
		config.When("GetIgnoreErr", "account").Return(configAccount)
	}
	tests := []showAccountTest{
		{ // 0
			Input:         "bytemark show account",
			ConfigAccount: "",
		},
		{ // 1
			Input:         "bytemark show account",
			ConfigAccount: "sec",
			AccountToGet:  "sec",
		},
		{ // 2
			Input:         "bytemark --account caan show account",
			ConfigAccount: "",
			AccountToGet:  "caan",
		},
		{ // 3
			Input:         "bytemark --account caan show account",
			ConfigAccount: "jast",
			AccountToGet:  "caan",
		},
		{ // 4
			Input:         "bytemark show account thay",
			ConfigAccount: "jast",
			AccountToGet:  "thay",
		},
		{ // 5
			Input:         "bytemark show account --account jast",
			ConfigAccount: "jast",
			AccountToGet:  "jast",
		},
	}

	runTest := func(i int, test showAccountTest) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("TestShowAccountCommand %d panicked: %v %s", i, err, debug.Stack())
			}
		}()
		config, c := baseTestAuthSetup(t, false)
		baseShowAccountSetup(c, config, test.ConfigAccount)
		if test.AccountToGet == "" {
			c.When("GetAccount", "").Return(&lib.Account{
				Name:      "default-account",
				BrainID:   112,
				BillingID: 213,
				Groups:    []*brain.Group{},
				Owner: &billing.Person{
					FirstName: "Doctor",
					LastName:  "Testo",
				},
				TechnicalContact: &billing.Person{
					FirstName: "Doctor",
					LastName:  "Testo",
				},
			}).Times(1)
		} else {
			c.When("GetAccount", test.AccountToGet).Return(&lib.Account{
				Name:      "default-account",
				BrainID:   112,
				BillingID: 213,
				Groups:    []*brain.Group{},
				Owner: &billing.Person{
					FirstName: "Doctor",
					LastName:  "Testo",
				},
				TechnicalContact: &billing.Person{
					FirstName: "Doctor",
					LastName:  "Testo",
				},
			}).Times(1)
		}
		err := global.App.Run(strings.Split(test.Input, " "))
		if !test.ShouldErr && err != nil {
			t.Errorf("TestShowAccountCommand %d shouldn't err, but did: %T{%s}", i, err, err.Error())
		} else if test.ShouldErr && err == nil {
			t.Errorf("TestShowAccountCommand %d should err, but didn't", i)
		}
		if ok, vErr := c.Verify(); !ok {
			t.Errorf("TestShowAccountCommand %d client failed to verify: %s", i, vErr.Error())
		}
	}

	for i, test := range tests {
		runTest(i, test)
	}

}

func TestShowGroupCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("GetGroup").Return(&defGroup)
	gpname := lib.GroupName{Group: "test-group", Account: "test-account"}

	group := getFixtureGroup()
	c.When("GetGroup", &gpname).Return(&group, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark show group test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowServerCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("GetVirtualMachine").Return(&defVM)
	vmname := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}
	vm := getFixtureVM()
	c.When("GetVirtualMachine", &vmname).Return(&vm, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark show server test-server.test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowVLANsCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	vlans := []brain.VLAN{getFixtureVLAN()}
	c.When("GetVLANs").Return(&vlans, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show vlans", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowVLANCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	vlanNum := 1
	vlan := getFixtureVLAN()
	c.When("GetVLAN", vlanNum).Return(&vlan, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show vlan 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowIPRangesCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	ipRanges := []brain.IPRange{getFixtureIPRange()}
	c.When("GetIPRanges").Return(&ipRanges, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show ip_ranges", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowIPRangeCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	ipRangeID := "1"
	ipRange := getFixtureIPRange()
	c.When("GetIPRange", ipRangeID).Return(&ipRange, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show ip_range 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowIPRangeWithIPRangeCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	rangeString := "192.168.33.0/24"
	ipRange := getFixtureIPRange()
	c.When("GetIPRange", rangeString).Return(&ipRange, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show ip_range 192.168.33.0/24", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowHeadsCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	heads := []brain.Head{getFixtureHead()}
	c.When("GetHeads").Return(&heads, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show heads", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowHeadCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	headID := "1"
	head := getFixtureHead()
	c.When("GetHead", headID).Return(&head, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show head 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowTailsCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	tails := []brain.Tail{getFixtureTail()}
	c.When("GetTails").Return(&tails, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show tails", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowTailCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	tailID := "1"
	tail := getFixtureTail()
	c.When("GetTail", tailID).Return(&tail, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show tail 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowStoragePoolsCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	storagePool := []brain.StoragePool{getFixtureStoragePool()}
	c.When("GetStoragePools").Return(&storagePool, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show storage_pools", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowStoragePoolCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	storagePoolID := "1"
	storagePool := getFixtureStoragePool()
	c.When("GetStoragePool", storagePoolID).Return(&storagePool, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show storage_pool 1", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowMigratingVMsCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	vms := []brain.VirtualMachine{getFixtureVM()}
	c.When("GetMigratingVMs").Return(&vms, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show migrating_vms", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowStoppedEligibleVMsCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	vms := []brain.VirtualMachine{getFixtureVM()}
	c.When("GetStoppedEligibleVMs").Return(&vms, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show stopped_eligible_vms", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowRecentVMsCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	vms := []brain.VirtualMachine{getFixtureVM()}
	c.When("GetRecentVMs").Return(&vms, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show recent_vms", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

// TODO(telyn): show account? show user?
func TestShowPrivileges(t *testing.T) {

	tests := []struct {
		privs     brain.Privileges
		args      string
		user      string
		shouldErr bool
	}{
		{
			privs:     brain.Privileges{{VirtualMachineID: 643, Username: "satan", Level: brain.VMConsolePrivilege}},
			user:      "",
			args:      "bytemark --admin show privileges",
			shouldErr: false,
		},
	}

	for i, test := range tests {
		_, c := baseTestAuthSetup(t, true)
		c.When("GetPrivileges", test.user).Return(test.privs, nil)

		err := global.App.Run(strings.Split(test.args, " "))
		if test.shouldErr && err == nil {
			t.Errorf("TestShowPrivilege %d should err and didn't", i)
		} else if !test.shouldErr && err != nil {
			t.Errorf("TestShowPrivilege %d shouldn't err, but: %s", i, err)
		}
		if ok, err := c.Verify(); !ok {
			t.Fatal(err)
		}
	}
}
