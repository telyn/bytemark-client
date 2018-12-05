package commands

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	billingRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/billing"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "assent",
		Usage:       "assent to Bytemark terms and conditions",
		UsageText:   "assent [--agreement <agreement id>] [--person <username>] [--account <account>|--account-id <account id>] [--name <full name>] [--email <email>]",
		Description: "Assent to Bytemark terms and conditions.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "agreement",
				Usage: "the agreement id to assent to",
			},
			cli.StringFlag{
				Name:  "person",
				Usage: "the username of the person who is assenting",
			},
			cli.GenericFlag{
				Name:  "account",
				Usage: "the account which is assenting",
				Value: new(flags.AccountNameFlag),
			},
			cli.IntFlag{
				Name:  "accountid",
				Usage: "the account id of the account which is assenting",
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
		Action: app.Action(with.RequiredFlags("agreement", "person"), func(ctx *app.Context) error {
			accountString := ctx.String("account")
			accountNameFlag := flags.AccountName(ctx, "account")
			accountID := ctx.Int("accountid")
			var account billing.Account

			if accountString == "" && accountID == 0 {
				return fmt.Errorf("Neither --account or --accountid was set (or should not be blank/zero)")
			} else if accountNameFlag.SetFromCommandLine && accountString != "" && accountID != 0 {
				return fmt.Errorf("--account and --accountid have both been set when only one is required")
			}

			err := with.Auth(ctx)
			if err != nil {
				return err
			}
			if accountString != "" && accountID == 0 {
				// cant use with.Account() because this gets the account details of the person currently signed in, even if staff
				account, err = billingRequests.GetAccountByBigVName(ctx.Client(), accountString)
				if err != nil {
					return err
				}
				accountID = account.ID
			}

			person, personErr := billingRequests.GetPerson(ctx.Client(), ctx.String("person"))
			if personErr != nil {
				return personErr
			}

			name := ctx.String("name")
			email := ctx.String("email")
			prompt := ""

			if name == "" {
				name = person.FirstName + " " + person.LastName
				prompt = fmt.Sprintf("Name was not specified. Name of person will be used: %s. Is this correct?", name)
			}

			if email == "" {
				email = person.Email
				if prompt != "" {
					prompt = fmt.Sprintf("No name or email was specified. Name and email of person will be used is: %s and %s. Is this correct?", name, email)
				} else {
					prompt = fmt.Sprintf("Email was not specified. Email of person will be used: %s. Is this correct?", email)
				}
			}

			if prompt != "" && !util.PromptYesNo(ctx.Prompter(), prompt) {
				ctx.LogErr("Exiting. Please explicitly state Name and Email using the --name and --email options")
				return util.UserRequestedExit{}
			}

			err = billingRequests.AssentToAgreement(ctx.Client(), billing.Assent{
				AgreementID: ctx.String("agreement"),
				AccountID:   accountID,
				PersonID:    person.ID,
				Name:        name,
				Email:       email,
			})

			if err == nil {
				ctx.LogErr("Successfully added assent for account %d", account.ID)
			}
			return err
		}),
	})
}
