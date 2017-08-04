package main

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "signup",
		Usage:     "sign up for Bytemark's hosting service",
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
			frm.SetMaxWidth(120)

			err = frm.Run()
			if err != nil {
				return err
			}

			if !*signup {
				return util.UserRequestedExit{}
			}

			if problems, ok := frm.Validate(); !ok {
				log.Log(strings.Join(problems, "\r\n"))
				return util.UserRequestedExit{}
			}

			// TODO(telyn): this whoole section should be moved into a function in util/form.go - CreateAPIObjectsFromSignupForm(*Form) (Account, CreditCard) or something.
			account := lib.Account{}

			account.Owner = billing.Person{
				Username:             fields[util.FormFieldOwnerName].Value(),
				Password:             fields[util.FormFieldOwnerPassword].Value(),
				Email:                fields[util.FormFieldOwnerEmail].Value(),
				FirstName:            fields[util.FormFieldOwnerFirstName].Value(),
				LastName:             fields[util.FormFieldOwnerLastName].Value(),
				Address:              fields[util.FormFieldOwnerAddress].Value(),
				City:                 fields[util.FormFieldOwnerCity].Value(),
				Postcode:             fields[util.FormFieldOwnerPostcode].Value(),
				Country:              fields[util.FormFieldOwnerCountryCode].Value(),
				Phone:                fields[util.FormFieldOwnerPhoneNumber].Value(),
				MobilePhone:          fields[util.FormFieldOwnerMobileNumber].Value(),
				Organization:         fields[util.FormFieldOwnerOrgName].Value(),
				OrganizationDivision: fields[util.FormFieldOwnerOrgDivision].Value(),
				VATNumber:            fields[util.FormFieldOwnerOrgVATNumber].Value(),
			}

			if creditCardForm {
				card := spp.CreditCard{
					Number: fields[util.FormFieldCreditCardNumber].Value(),
					Name:   fields[util.FormFieldCreditCardName].Value(),
					Expiry: fields[util.FormFieldCreditCardExpiry].Value(),
					CVV:    fields[util.FormFieldCreditCardCVV].Value(),
				}

				token, err := global.Client.GetSPPToken(card, account.Owner)
				if err != nil {
					return err
				}

				cardRef, err = global.Client.CreateCreditCardWithToken(card, token)
				if err != nil {
					return err
				}
			}
			account.CardReference = cardRef
			createdAccount, err := global.Client.RegisterNewAccount(account)

			if _, ok := err.(lib.AccountCreationDeferredError); ok {
				log.Log(err.Error())
				return nil
			} else if err != nil {
				log.Log("Couldn't create an account for you")
				return err
			}
			log.Logf("Account created successfully - you'll now be able to log in as '%s' and set up some servers! You should also be receiving a welcome email shortly.\r\n", createdAccount.Owner.Username)
			return nil

		}),
	})
}
