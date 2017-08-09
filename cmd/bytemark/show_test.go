package main

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/cheekybits/is"
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

// mkExpectedOutput loops over every value in template, running fmt.Sprintf to replace all %s with substitute
// in this way it can be used to prepare an expected output array. See showTestAccountOutput and TestShowAccountCommand for some overly complicated example usage.
func mkExpectedOutput(template map[string]string, substitute string) map[string]string {
	m := make(map[string]string)
	for k, v := range template {
		toReplace := strings.Count(v, "%s")
		fmtArgs := make([]interface{}, toReplace)
		for i := 0; i < toReplace; i++ {
			fmtArgs[i] = substitute
		}
		m[k] = fmt.Sprintf(v, fmtArgs...)
	}
	return m
}

var showAccountTestOutput = map[string]string{
	"table": "+-----------+------+-----------+--------+\n| BillingID | Name | Suspended | Groups |\n+-----------+------+-----------+--------+\n|       213 | %s | false     |        |\n+-----------+------+-----------+--------+\n",
	"json":  "{\n    \"name\": \"%s\",\n    \"owner\": {\n        \"username\": \"\",\n        \"email\": \"\",\n        \"password\": \"\",\n        \"firstname\": \"%s\",\n        \"surname\": \"Testo\",\n        \"address\": \"\",\n        \"city\": \"\",\n        \"postcode\": \"\",\n        \"country\": \"\",\n        \"phone\": \"\"\n    },\n    \"technical_contact\": {\n        \"username\": \"\",\n        \"email\": \"\",\n        \"password\": \"\",\n        \"firstname\": \"%s\",\n        \"surname\": \"Testo\",\n        \"address\": \"\",\n        \"city\": \"\",\n        \"postcode\": \"\",\n        \"country\": \"\",\n        \"phone\": \"\"\n    },\n    \"billing_id\": 213,\n    \"brain_id\": 112,\n    \"card_reference\": \"\",\n    \"groups\": [],\n    \"suspended\": false\n}\n",
	"human": "213 - %s\n",
}

func TestShowAccountCommand(t *testing.T) {
	baseShowAccountSetup := func(c *mocks.Client, config *mocks.Config, configAccount, outputFormat string) {
		config.Mock.Functions = resetOneMockedFunction(config.Mock.Functions, "GetV", "output-format")
		config.When("GetV", "output-format").Return(util.ConfigVar{"output-format", outputFormat, "FLAG output-format"})
		config.When("GetIgnoreErr", "account").Return(configAccount)
	}

	tests := []showAccountTest{
		{ // 0
			Input:         "show account",
			ConfigAccount: "",
		},
		{ // 1
			Input:         "show account",
			ConfigAccount: "secg",
			AccountToGet:  "secg",
		},
		{ // 2
			Input:         "--account caan show account",
			ConfigAccount: "",
			AccountToGet:  "caan",
		},
		{ // 3
			Input:         "--account caan show account",
			ConfigAccount: "jast",
			AccountToGet:  "caan",
		},
		{ // 4
			Input:         "show account thay",
			ConfigAccount: "jast",
			AccountToGet:  "thay",
		},
		{ // 5
			Input:         "show account --account caan",
			ConfigAccount: "jast",
			AccountToGet:  "caan",
		},
	}

	runTest := func(i int, test showAccountTest) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("TestShowAccountCommand %d panicked: %v %s", i, err, debug.Stack())
			}
		}()

		expectedOutput := mkExpectedOutput(showAccountTestOutput, test.AccountToGet)
		if test.AccountToGet == "" {
			expectedOutput = mkExpectedOutput(showAccountTestOutput, "defa")
		}
		for _, format := range []string{"table", "json", "human"} {

			config, c := baseTestAuthSetup(t, false)
			baseShowAccountSetup(c, config, test.ConfigAccount, format)
			if test.AccountToGet == "" {
				c.When("GetAccount", "").Return(lib.Account{
					Name:      "defa",
					BrainID:   112,
					BillingID: 213,
					Groups:    []brain.Group{},
					Owner: billing.Person{
						FirstName: "defa",
						LastName:  "Testo",
					},
					TechnicalContact: billing.Person{
						FirstName: "defa",
						LastName:  "Testo",
					},
				}).Times(1)
			} else {
				// wat
				c.When("GetAccount", test.AccountToGet).Return(lib.Account{
					Name:      test.AccountToGet,
					BrainID:   112,
					BillingID: 213,
					Groups:    []brain.Group{},
					Owner: billing.Person{
						FirstName: test.AccountToGet,
						LastName:  "Testo",
					},
					TechnicalContact: billing.Person{
						FirstName: test.AccountToGet,
						LastName:  "Testo",
					},
				}).Times(1)
			}

			buf := bytes.Buffer{}
			global.App.Writer = &buf
			global.App.ErrWriter = &buf
			args := fmt.Sprintf("bytemark --output-format=%s %s", format, test.Input)
			//t.Logf("TestShowAccountCommand %d args: %s", i, args)
			err := global.App.Run(strings.Split(args, " "))
			if !test.ShouldErr && err != nil {
				t.Errorf("TestShowAccountCommand %d shouldn't err, but did: %T{%s}", i, err, err.Error())
			} else if test.ShouldErr && err == nil {
				t.Errorf("TestShowAccountCommand %d should err, but didn't", i)
			}
			if ok, vErr := c.Verify(); !ok {
				t.Errorf("TestShowAccountCommand %d client failed to verify: %s", i, vErr.Error())
			}
			if expected, ok := expectedOutput[format]; ok {
				actual := buf.String()
				if expected != actual {
					t.Errorf("TestShowAccountCommand %d didn't match expected %s output.\r\nExpected: %q\r\nActual %q", i, format, expected, actual)
				}
			} else {
				t.Errorf("TestShowAccountCommand %d didn't have an expected %s output. Maybe it should be %q", i, format, buf.String())
			}
		}
	}

	for i, test := range tests {
		runTest(i, test)
	}

}

func TestShowGroupCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("GetGroup").Return(defGroup)
	gpname := lib.GroupName{Group: "test-group", Account: "test-account"}

	group := getFixtureGroup()
	c.When("GetGroup", gpname).Return(&group, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark show group test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestShowServerCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	config.When("GetVirtualMachine").Return(defVM)
	vmname := lib.VirtualMachineName{VirtualMachine: "test-server", Group: "test-group", Account: "test-account"}
	vm := getFixtureVM()
	c.When("GetVirtualMachine", vmname).Return(&vm, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark show server test-server.test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowVLANsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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

func TestAdminShowMigratingDiscsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	discs := brain.Discs{{ID: 134, StorageGrade: "sata", Size: 25600}}
	c.When("GetMigratingDiscs").Return(discs, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show migrating_discs", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestAdminShowStoppedEligibleVMsCommand(t *testing.T) {
	// TODO(telyn): make table-driven
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
	// TODO(telyn): make table-driven
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

func TestAdminShowDiscByIDCommand(t *testing.T) {
	is := is.New(t)
	_, c := baseTestAuthSetup(t, true)

	discID := 132
	disc := getFixtureDisc()
	c.When("GetDiscByID", discID).Return(&disc, nil).Times(1)

	err := global.App.Run(strings.Split("bytemark --admin show disc_by_id 132", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
