package add

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:    "api key",
		Aliases: []string{"apikey"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "expires-at",
				Usage: "Date the API key should expire. Leave unset for keys that never expire",
			},
			cli.StringSliceFlag{
				Name:  "group",
				Usage: "Group to grant the API key administrative privilege over",
			},
			cli.StringFlag{
				Name:  "label",
				Usage: "user-friendly label for the API key",
			},
			cli.StringSliceFlag{
				Name:  "server",
				Usage: "Server to grant the API key administrative privilege over",
			},
			cli.StringFlag{
				Name:  "user",
				Usage: "User the API key will be attached to. Defaults to the user you log in as",
			},
		},
		Usage:     "add an API key to your Bytemark Cloud Servers user",
		UsageText: "add api key [--server <cloud server>]... [--group <group name>]... [--user <user>] <label>",
		Description: `--expires-at may be set to any date format the Brain
accepts, but we generally recommend ISO8601 format.

Servers and groups will be searched for on the default account for the user
you are logged in as. This may trip up cluster administrators, so
bytemark-client will refuse to create an API key whose access is not a subset
of the access the specified user normally has. To create such an API key you
can either add the necessary privileges to ensure that the API key privileges
are a subset, or create the API key without privileges and add them via the
grant command, which does not have this limitation.

Note that the API key will only currently be able to access the Bytemark Cloud
Servers API - to manage servers and groups.

Multiple --group and --server flags (and combinations thereof) can be supplied,
and the API key will be have privileges over each that is supplied.`,
		Action: app.Action(args.Optional("label"), with.Auth, func(ctx *app.Context) error {
			servers, serverIDsMap, err := findServers(ctx)
			if err != nil {
				return err
			}
			groups, groupIDsMap := findGroups(ctx)
			if err != nil {
				return err
			}

			apiKey, err := brainRequests.CreateAPIKey(ctx.Client(),keySpec)
			if err != nil {
				return err
			}
			privs := makeServerPrivs(servers, apiKey)
			privs = append(privs, makeGroupPrivs(groups, apiKey))
			privErrs := []error{}
			for _, priv := range privs {
				err = ctx.Client().GrantPrivilege(privs)
				if err != nil {
					privErrs
			}
		}),
	})
}
