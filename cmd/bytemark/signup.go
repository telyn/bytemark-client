package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"github.com/codegangsta/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "signup",
		Usage:     "Sign up for Bytemark's hosting service.",
		UsageText: "bytemark signup",
		Description: `This will create a new SSO and billing account and set your credit card details.

If you are creating an account on behalf of an organisation needing a different payment method, you'll need to email Bytemark support instead.

If you have previously used the client, you'll have a login and will need to add the --force flag in order to create a new account`,
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

			fields, frm := util.MakeSignupForm()

			frm.Run()

			account := lib.Account{
				Name: fields[util.FIELD_ACCOUNT].Value(),
			}

			account.Owner = &lib.Person{
				Username:             fields[util.FIELD_OWNER_NAME].Value(),
				Email:                fields[util.FIELD_OWNER_EMAIL].Value(),
				FirstName:            fields[util.FIELD_OWNER_FIRSTNAME].Value(),
				LastName:             fields[util.FIELD_OWNER_LASTNAME].Value(),
				Country:              fields[util.FIELD_OWNER_CC].Value(),
				City:                 fields[util.FIELD_OWNER_POSTCODE].Value(),
				Address:              fields[util.FIELD_OWNER_ADDRESS].Value(),
				Phone:                fields[util.FIELD_OWNER_PHONE].Value(),
				MobilePhone:          fields[util.FIELD_OWNER_MOBILE].Value(),
				Organization:         fields[util.FIELD_OWNER_ORG_NAME].Value(),
				OrganizationDivision: fields[util.FIELD_OWNER_ORG_DIVISION].Value(),
				VATNumber:            fields[util.FIELD_OWNER_ORG_VAT].Value(),
			}
			account.TechnicalContact = account.Owner

			card := lib.CreditCard{
				Number: fields[util.FIELD_CC_NUMBER].Value(),
				Name:   fields[util.FIELD_CC_NAME].Value(),
				Expiry: fields[util.FIELD_CC_EXPIRY].Value(),
				CVV:    fields[util.FIELD_CC_CVV].Value(),
			}

			ref, err := global.Client.CreateCreditCard(&card)
			if err != nil {
				return err
			}
			account.CardReference = ref
			createdAccount, err := global.Client.RegisterNewAccount(&account)
			return err
		}),
	})
}
