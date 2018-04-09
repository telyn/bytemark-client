package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:      "config",
		Usage:     "manage the bytemark client's configuration",
		UsageText: "show config",
		// FIXME: List of variables also in description of
		//        update config - need to DRY
		Description: `View the bytemark-client configuration.

    The following variables are displayed:
        account - the default account, used when you do not explicitly state an account - defaults to the same as your user name
        token - the token used for authentication
        user - the user that you log in as by default
        group - the default group, used when you do not explicitly state a group (defaults to 'default')

        debug-level - the default debug level. Set to 0 unless you like lots of output.
	api-endpoint - the endpoint for domains (among other things?)
        auth-endpoint - the endpoint to authenticate to. https://auth.bytemark.co.uk is the default.
        endpoint - the brain endpoint to connect to. https://uk0.bigv.io is the default.
        billing-endpoint - the billing API endpoint to connect to. https://bmbilling.bytemark.co.uk is the default.
        spp-endpoint - the SPP endpoint to use. https://spp-submissions.bytemark.co.uk is the default.`,
		Flags: app.OutputFlags("vars", "array"),
		Action: app.Action(viewConfig),
	})
}

func viewConfig(c *app.Context) error {
	vars, err := c.Config().GetAll()
	if err != nil {
		return err
	}
	return c.OutputInDesiredForm(vars, output.List)
}
