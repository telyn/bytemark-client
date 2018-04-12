package admin_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
)

func TestApproveVM(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm123", Group: "group", Account: "account"}
	c.When("ApproveVM", vmName, false).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "approve", "vm", "vm123.group.account"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestApproveVMAndPowerOn(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	vmName := lib.VirtualMachineName{VirtualMachine: "vm122", Group: "group", Account: "account"}
	c.When("ApproveVM", vmName, true).Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "approve", "vm", "vm122.group.account", "true"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestApproveVMError(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)

	config.When("GetVirtualMachine").Return(defVM)

	approveErr := fmt.Errorf("Error approving")
	vmName := lib.VirtualMachineName{VirtualMachine: "vm121", Group: "group", Account: "account"}
	c.When("ApproveVM", vmName, false).Return(approveErr).Times(1)

	err := app.Run([]string{"bytemark", "approve", "vm", "vm121.group.account"})

	is.Equal(err, approveErr)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
