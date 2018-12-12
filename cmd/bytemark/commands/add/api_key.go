package add

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output/morestrings"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
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
			cli.GenericFlag{
				Name:  "account-admin",
				Usage: "Account to grant the API key administrative privilege over",
				Value: &flags.GroupNameSliceFlag{},
			},
			cli.GenericFlag{
				Name:  "group",
				Usage: "Group to grant the API key administrative privilege over",
				Value: &flags.GroupNameSliceFlag{},
			},
			cli.StringFlag{
				Name:  "label",
				Usage: "user-friendly label for the API key",
			},
			cli.GenericFlag{
				Name:  "server",
				Usage: "Server to grant the API key administrative privilege over",
				Value: &flags.VirtualMachineNameSliceFlag{},
			},
			cli.StringFlag{
				Name:  "user",
				Usage: "User the API key will be attached to. Defaults to the user you log in as",
			},
		},
		Usage:     "add an API key to your Bytemark Cloud Servers user",
		UsageText: "add api key [--server <cloud server>]... [--group <group name>]... [--account-admin <account name>]... [--user <user>] <label>",
		Description: `--expires-at may be set to any date format the Brain
accepts, but we generally recommend ISO8601 format.

--label or <label> is just for your reference to make it easier to discern keys
from each other at a glance, without having to compare lists of privileges.

Servers and groups will be searched for on the default account for the user
you are logged in as - not the user specified in --user. This may trip up
cluster administrators - use impersonation if you're not granting account_admin
access.

Only cluster administrators can create API keys with account-admin support. Pop
and email over in to Bytemark support and we'll do it for you.

Multiple --account-admin, --group and --server flags (and combinations thereof)
can be supplied, and the API key will have privileges over each that is
supplied.

Note that the API key will only be able to access the Bytemark Cloud Servers API
- to manage your cloud servers, IPs, discs and backups, not billing or domains.

In future we may implement API key support on more of our API, but the 
privileges to access those parts of the API will not be added automatically to
any existing API keys. In other words - an API key with account admin access now
will not be able to suddenly access your invoices/account owner address in
future if and when we add that feature to API keys.

EXAMPLES

To create an API key for yourself without any privileges:

  bytemark add api key currently-useless-api-key

You can always add privileges later to make it a non-useless API key with the
'grant' command using the --api-key flag. See 'bytemark help grant' for more
details.

To create an API key for yourself which can create, delete, and alter any server
in the 'kube' and 'swarm' groups on your default account:

  bytemark add api key --group kube --group swarm container-swarms-key

To create a key for yourself with access to rescale the database servers in the
'internal' group on an account 'big-data-services' which is not your default
account:

  bytemark add api key --server db1.internal.big-data-services \
                       --server db2.internal.big-data-services \
					   --label auto-rescale-db
`,
		Action: app.Action(args.Optional("label"), with.User("user"), func(ctx *app.Context) error {
			keySpec, err := makeAPIKeySpec(ctx)
			if err != nil {
				ctx.LogErr("Couldn't make a specification for the API key")
				return err
			}

			apiKey, err := brainRequests.CreateAPIKey(ctx.Client(), "", keySpec)
			if err != nil {
				return err
			}
			keySpec.APIKey = apiKey.APIKey
			keySpec.ID = apiKey.ID

			for i := range keySpec.Privileges {
				keySpec.Privileges[i].APIKeyID = apiKey.ID
			}

			if len(keySpec.Privileges) == 0 {
				ctx.LogErr("Successfully created an api key:")
				err = apiKey.PrettyPrint(ctx.Writer(), prettyprint.Full)
				return err
			}

			ctx.LogErr("Successfully created api key, now creating %d privileges...", len(keySpec.Privileges))

			// keep this error separate cause we check it later
			privsErr := addAPIKeyPrivileges(ctx, keySpec)

			done := "done."
			if privsErr != nil {
				done = "done, with some errors. Here's the API key we created with all the privileges we managed to create:"
			}

			ctx.LogErr(done)
			apiKey.PrettyPrint(ctx.Writer(), prettyprint.Full)

			if privsErr != nil {
				return privsErr
			}

			return nil

		}),
	})
}

func makeAPIKeySpec(ctx *app.Context) (spec brain.APIKey, err error) {
	spec.Label = ctx.String("label")
	spec.UserID = ctx.User.ID
	spec.ExpiresAt = ctx.String("expires-at")

	if spec.Label == "" {
		return spec, errors.New("a label must be specified for the key")
	}

	typesWithErrs := []string{}
	spec.Privileges, err = addAPIKeySpecifyServerPrivileges(ctx, spec.Privileges)
	if err != nil {
		typesWithErrs = append(typesWithErrs, "servers")
	}
	spec.Privileges, err = addAPIKeySpecifyGroupPrivileges(ctx, spec.Privileges)
	if err != nil {
		typesWithErrs = append(typesWithErrs, "groups")
	}
	spec.Privileges, err = addAPIKeySpecifyAccountPrivileges(ctx, spec.Privileges)
	if err != nil {
		typesWithErrs = append(typesWithErrs, "accounts")
	}
	if len(typesWithErrs) > 0 {
		err = fmt.Errorf("Some %s could not be looked up", morestrings.JoinWithSpecialLast(", ", " and ", typesWithErrs))
	}

	return
}

// creates a privilege spec for every --server flag
func addAPIKeySpecifyServerPrivileges(ctx *app.Context, privileges brain.Privileges) (brain.Privileges, error) {
	var serverErr error
	for _, serverFlag := range flags.VirtualMachineNameSlice(ctx, "server") {
		serverName := serverFlag.VirtualMachineName
		ctx.LogErr("Looking up %s", serverName)

		server, err := ctx.Client().GetVirtualMachine(serverName)
		if err != nil {
			ctx.LogErr(err.Error())
			if serverErr == nil {
				serverErr = errors.New("Some servers could not be looked up - see above for errors")
			}
			continue
		}

		privileges = append(privileges, brain.Privilege{
			Username:         ctx.User.Username,
			VirtualMachineID: server.ID,
		})
	}
	return privileges, nil
}

// creates a privilege spec for every --group flag
func addAPIKeySpecifyGroupPrivileges(ctx *app.Context, privileges brain.Privileges) (brain.Privileges, error) {
	var groupErr error
	for _, groupFlag := range flags.GroupNameSlice(ctx, "group") {
		groupName := groupFlag.GroupName
		ctx.LogErr("Looking up %s", groupName)

		group, err := ctx.Client().GetGroup(groupName)
		if err != nil {
			ctx.LogErr(err.Error())
			if groupErr == nil {
				groupErr = errors.New("Some groups could not be looked up - see above for errors")
			}
			continue
		}

		privileges = append(privileges, brain.Privilege{
			Username: ctx.User.Username,
			GroupID:  group.ID,
		})
	}
	return privileges, nil
}

// creates a privilege spec for every --account-admin flag
func addAPIKeySpecifyAccountPrivileges(ctx *app.Context, privileges brain.Privileges) (brain.Privileges, error) {
	var accountErr error
	for _, accountFlag := range flags.AccountNameSlice(ctx, "account-admin") {
		accountName := accountFlag.AccountName
		ctx.LogErr("Looking up %s", accountName)

		account, err := ctx.Client().GetAccount(accountName)
		if err != nil {
			ctx.LogErr(err.Error())
			if accountErr == nil {
				accountErr = errors.New("Some accounts could not be looked up - see above for errors")
			}
			continue
		}

		privileges = append(privileges, brain.Privilege{
			Username:  ctx.User.Username,
			AccountID: account.BrainID,
		})
	}
	return privileges, nil
}

// addAPIKeyPrivileges goes over the privilege specs and grants them to the api
// key.
func addAPIKeyPrivileges(ctx *app.Context, apiKey brain.APIKey) error {
	// collect up all the errors to output a nice list at the end, so
	// that we still output the complete API key as it actually stands
	succeededPrivs := make(brain.Privileges, 0, len(apiKey.Privileges))
	privErrs := []addAPIKeyPrivErr{}

	for i, priv := range apiKey.Privileges {
		priv.APIKeyID = apiKey.ID

		err := ctx.Client().GrantPrivilege(priv)
		if err != nil {
			privErrs = append(privErrs, addAPIKeyPrivErr{
				idx:  i,
				priv: priv,
				err:  err,
			})
		} else {
			succeededPrivs = append(succeededPrivs, priv)
		}
	}
	apiKey.Privileges = succeededPrivs

	if len(privErrs) > 0 {
		// catalogue the failures
		lines := make([]string, len(privErrs))
		for i, privErr := range privErrs {
			lines[i] = fmt.Sprintf("  • %s: %s", privErr.priv, privErr.err)
		}
		intro := fmt.Sprintf("Couldn't create %d/%d privileges requested:", len(privErrs), len(apiKey.Privileges))

		return errors.New(intro + "\n" + strings.Join(lines, "\n"))
	}
	return nil
}

type addAPIKeyPrivErr struct {
	idx  int
	priv brain.Privilege
	err  error
}
