package flags_test

import (
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestGroupNameSliceFlag(t *testing.T) {
	sf := flags.GroupNameSliceFlag{}
	err := sf.Set("staples.stapler")
	if err != nil {
		t.Errorf("got error from Set(): %s", err)
	}
	if len(sf) != 1 {
		t.Errorf("Expected len(GroupNameSliceFLag) to be 1, got %d",
			len(sf))
	}

	t.Logf("Value: %s", sf[0].Value)
	// it's a Preprocesser so we need to call Preprocess before we can validate
	// String()
	cfg, client, cliApp := testutil.BaseTestSetup(t, false, []cli.Command{})
	cfg.When("GetIgnoreErr", "account").Return("default-account")
	cfg.When("GetGroup").Return(lib.GroupName{Group: "default-group", Account: "default-account"})
	cfg.When("GetVirtualMachine").Return(lib.VirtualMachineName{VirtualMachine: "default-server", Group: "default-group", Account: "default-account"})

	// now some boilerplate to get a context
	// TODO(telyn): this should probably be refactored out since it'll be
	// wanted for basically every Preprocesser flag)
	client.When("AuthWithToken", "test-token").Return(nil)
	cliCtx := mocks.CliContext{}
	cliCtx.When("App").Return(cliApp)
	ctx := app.Context{
		Context: &cliCtx,
	}

	// with a context we may now Preprocess
	err = sf.Preprocess(&ctx)
	if err != nil {
		t.Errorf("Preprocess errored: %s", err)
	}

	if sf.String() != "staples.stapler" {
		t.Errorf("Expected %q, got %q", "staples.stapler", sf.String())
	}

}
