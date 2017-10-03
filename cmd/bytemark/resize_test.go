package main

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func TestResizeDisk(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	config.When("Force").Return(true)

	disc := brain.Disc{
		Size:         25600,
		StorageGrade: "sata",
	}

	config.When("GetVirtualMachine").Return(defVM)

	name := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}
	c.When("GetDisc", name, "disc-label").Return(&disc).Times(1)

	c.When("ResizeDisc", name, "disc-label", 35*1024).Return(nil).Times(1)

	err := app.Run(strings.Split("bytemark resize disc --force test-server disc-label 35", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
