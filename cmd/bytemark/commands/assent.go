package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "assent",
		Usage:       "assent to Bytemark terms and conditions",
		UsageText:   "bytemark assent --agreement <agreement id> --account <account> --person <username> [--name <full name> --email <email>]",
		Description: ``,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "agreement",
				Usage: "the agreement id to assent to",
			},
			cli.GenericFlag{
				Name:  "account",
				Usage: "The account which is assenting",
				Value: new(app.AccountNameFlag),
			},
			cli.StringFlag{
				Name:  "person",
				Usage: "the username of the person who is assenting",
			},
			cli.StringFlag{
				Name:  "name",
				Usage: "the full name of the person who is assenting. defaults to the full name of the person specified by the person flag",
			},
			cli.StringFlag{
				Name:  "email",
				Usage: "the email address of the person who is assenting. defaults to the full name of the person specified by the person flag",
			},
		},
		Action: app.Action(with.RequiredFlags("agreement", "account", "person"), with.Account("account"), func(ctx *app.Context) error {
			person, err := billingRequests.GetPerson(ctx.Client(), ctx.String("person"))
			if err != nil {
				return err
			}
			name := ctx.String("name")
			email := ctx.String("email")

			if name == "" {
				name = person.FirstName + " " + person.LastName
			}

			if email == "" {
				email = person.Email
			}

			err = billingRequests.AssentToAgreement(ctx.Client(), billing.Assent{
				AgreementID: ctx.String("agreement"),
				AccountID:   ctx.Account.ID,
				PersonID:    person.ID,
				Name:        name,
				Email:       email,
			})

			if err == nil {
				ctx.LogErr("Successfully added assent for account %d", ctx.Account.ID)
			}
			return err
		}),
	})
}
