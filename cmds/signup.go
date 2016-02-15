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

	log.Log("Okay, we're going to have to ask you a lot of questions in three parts - first about the account, second about you, and thirdly your card details.")
	log.Log("If you wish to pay by BACS, PayPal, etc. you can either set up an account with a card now and email support, or email Bytemark at sales@bytemark.co.uk")
	log.Log("At any time you may press Ctrl+C to exit.")

	log.Log("Now some details about you / your organisation.")
	questions := 15
	curq := 1
	account := lib.Account{
		Name: util.Promptf("[%d/%d] Account name (this will be used in your default DNS entries):", curq, questions),
	}

	account.Owner = &lib.Person{
		Username:    util.Promptf("[%d/%d] Username (this is the name you will use to log in):", curq, questions),
		Email:       util.Promptf("[%d/%d] Email address (it is important that this is correct):", curq, questions),
		FirstName:   util.Promptf("[%d/%d] First name:", curq, questions),
		LastName:    util.Promptf("[%d/%d] Last name:", curq, questions),
		Country:     util.Promptf("[%d/%d] Post code:", curq, questions),
		City:        util.Promptf("[%d/%d] City:", curq, questions),
		Address:     util.Promptf("[%d/%d] Street address:", curq, questions),
		Phone:       util.Promptf("[%d/%d] Phone number:", curq, questions),
		MobilePhone: util.Promptf("[%d/%d] Mobile phone number (if different):", curq, questions),
	}
	account.TechnicalContact = account.Owner
	if util.PromptYesNo("[%d/%d] Are you creating this account for a company?") {
		account.Owner.Organization = util.Promptf("[%d/%d] Organization name:", curq, questions)
		account.Owner.OrganizationDivision = util.Promptf("[%d/%d] Orgnization division:", curq, questions)
		account.Owner.VATNumber = util.Promptf("[%d/%d] VAT number:", curq, questions)
	}

	card := lib.CreditCard{
		Number: util.Promptf("[%d/%d] Credit card number:"),
		Name:   util.Promptf("[%d/%d] Name on the card:"),
		Expiry: util.Promptf("[%d/%d] Expiry date (MMYY):"),
		CVV:    util.Promptf("[%d/%d] CVV (3-4 digit number on the back):"),
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
