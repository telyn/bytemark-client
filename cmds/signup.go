package cmds

import (
	"bytemark.co.uk/client/cmds/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
)

func (c *CommandSet) HelpForSignup() util.ExitCode {
	log.Log("usage: bytemark signup [--force]")
	log.Log("")
	log.Log("Sign up for Bytemark's hosting service.")
	log.Log("This will create a new SSO and billing account and set your credit card details.")
	log.Log("")
	log.Log("If you're creating an account on behalf of an organisation needing a different")
	log.Log("payment method, you'll need to email Bytemark support instead.")
	log.Log("")
	log.Log("If you have previously used the client, you'll have a login")
	log.Log("and you will need to add the --force flag to.")
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) Signup(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	// TODO(telyn): check a terminal is attached to stdin to try to help prevent fraudy/spammy crap just in case
	ssoExists := false
	token := cmds.config.GetIgnoreErr("token")
	if token != "" {
		ssoExists = true
	}
	user, err := cmds.config.GetV("user")
	if err == nil && user.Source != "ENV USER" {
		ssoExists = true
	}

	fields, frm := util.MakeSignupForm()

	frm.Run()

	account := lib.Account{
		Name: fields[util.FIELD_ACCOUNT].Value(),
	}

	account.Owner = &lib.Person{
		Username:     fields[util.FIELD_OWNER_NAME].Value(),
		Email:        fields[util.FIELD_OWNER_EMAIL].Value(),
		FirstName:    fields[util.FIELD_OWNER_FIRSTNAME].Value(),
		LastName:     fields[util.FIELD_OWNER_LASTNAME].Value(),
		Country:      fields[util.FIELD_OWNER_CC].Value(),
		City:         fields[util.FIELD_OWNER_POSTCODE].Value(),
		Address:      fields[util.FIELD_OWNER_ADDRESS].Value(),
		Phone:        fields[util.FIELD_OWNER_PHONE].Value(),
		MobilePhone:  fields[util.FIELD_OWNER_MOBILE].Value(),
		Organization: fields[util.FIELD_OWNER_ORG_NAME].Value(),
		Division:     fields[util.FIELD_OWNER_ORG_DIVISION].Value(),
		VATNumber:    fields[util.FIELD_OWNER_ORG_VAT].Value(),
	}
	account.TechnicalContact = account.Owner

	card := lib.CreditCard{
		Number: fields[util.FIELD_CC_NUMBER].Value(),
		Name:   fields[util.FIELD_CC_NAME].Value(),
		Expiry: fields[util.FIELD_CC_EXPIRY].Value(),
		CVV:    fields[util.FIELD_CC_CVV].Value(),
	}

	ref, err := cmds.client.CreateCreditCard(&card)
	account.CardReference = ref
	createdAccount, err := cmds.client.CreateAccount(&account)
	if err != nil {
		// undo card creation?
		return util.ProcessError(err)
	}

	return util.E_SUCCESS
}
