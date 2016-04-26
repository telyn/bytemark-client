package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"github.com/codegangsta/cli"
	"strings"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "signup",
		Usage:     "Sign up for Bytemark's hosting service.",
		UsageText: "bytemark signup",
		Description: `This will create a new SSO and billing account and set your credit card details.

If you are creating an account on behalf of an organisation needing a different payment method, you'll need to email Bytemark support instead.

If you have previously used the client, you'll have a login and will need to add the --force flag in order to create a new account`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "card",
				Usage: "card reference string to use. If not specified you will be prompted for card details",
			},
			cli.BoolFlag{
				Name:  "force",
				Usage: "sign up for a new account & login despite already having a login.",
			},
		},
		Action: With(func(c *Context) error {

			// TODO(telyn): check a terminal is attached to stdin to try to help prevent fraudy/spammy crap just in case
			ssoExists := false
			token := global.Config.GetIgnoreErr("token")
			if token != "" {
				ssoExists = true
			}
			user, err := global.Config.GetV("user")
			if err == nil && user.Source != "ENV USER" {
				ssoExists = true
			}

			if ssoExists && !c.Bool("force") {
				return c.Help("You already have a login configured, you may wish to use 'create account' to add another account to your user, or add the force flag.")
			}
			cardRef := c.String("card")
			creditCardForm := true
			if cardRef != "" {
				creditCardForm = false
			}

			fields, frm, signup := util.MakeSignupForm(creditCardForm)

			err = frm.Run()
			if err != nil {
				return err
			}

			if !*signup {
				return new(util.UserRequestedExit)
			}

			if problems, ok := frm.Validate(); !ok {
				log.Log(strings.Join(problems, "\r\n"))
				return new(util.UserRequestedExit)
			}

			account := lib.Account{}

			account.Owner = &lib.Person{
				Username:             fields[util.FIELD_OWNER_NAME].Value(),
				Password:             fields[util.FIELD_OWNER_PASS].Value(),
				Email:                fields[util.FIELD_OWNER_EMAIL].Value(),
				FirstName:            fields[util.FIELD_OWNER_FIRSTNAME].Value(),
				LastName:             fields[util.FIELD_OWNER_LASTNAME].Value(),
				Country:              fields[util.FIELD_OWNER_CC].Value(),
				City:                 fields[util.FIELD_OWNER_CITY].Value(),
				Postcode:             fields[util.FIELD_OWNER_POSTCODE].Value(),
				Address:              fields[util.FIELD_OWNER_ADDRESS].Value(),
				Phone:                fields[util.FIELD_OWNER_PHONE].Value(),
				MobilePhone:          fields[util.FIELD_OWNER_MOBILE].Value(),
				Organization:         fields[util.FIELD_OWNER_ORG_NAME].Value(),
				OrganizationDivision: fields[util.FIELD_OWNER_ORG_DIVISION].Value(),
				VATNumber:            fields[util.FIELD_OWNER_ORG_VAT].Value(),
			}

			if creditCardForm {
				card := lib.CreditCard{
					Number: fields[util.FIELD_CC_NUMBER].Value(),
					Name:   fields[util.FIELD_CC_NAME].Value(),
					Expiry: fields[util.FIELD_CC_EXPIRY].Value(),
					CVV:    fields[util.FIELD_CC_CVV].Value(),
				}
				cardRef, err = global.Client.CreateCreditCard(&card)
				if err != nil {
					return err
				}
			}
			account.CardReference = cardRef
			createdAccount, err := global.Client.RegisterNewAccount(&account)

			if _, ok := err.(*lib.AccountCreationDeferredError); ok {
				log.Log(err.Error())
				return nil
			} else if err != nil {
				log.Log("Couldn't create an account for you")
				return err
			}
			log.Log("Account created successfully - you'll now be able to log in as '%s' and set up some servers! You should also be receiving a welcome email shortly.", createdAccount.Owner.Username)
			return nil

		}),
	})
}
