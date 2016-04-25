package util

import (
	"bytemark.co.uk/client/util/log"
	"github.com/telyn/form"
)

const (
	FIELD_ACCOUNT = iota
	FIELD_OWNER_NAME
	FIELD_OWNER_PASS
	FIELD_OWNER_PASS_CONFIRM
	FIELD_OWNER_EMAIL
	FIELD_OWNER_FIRSTNAME
	FIELD_OWNER_LASTNAME
	FIELD_OWNER_CC
	FIELD_OWNER_POSTCODE
	FIELD_OWNER_CITY
	FIELD_OWNER_ADDRESS
	FIELD_OWNER_PHONE
	FIELD_OWNER_MOBILE
	FIELD_OWNER_ORG_NAME
	FIELD_OWNER_ORG_DIVISION
	FIELD_OWNER_ORG_VAT
	FIELD_CC_NUMBER
	FIELD_CC_NAME
	FIELD_CC_EXPIRY
	FIELD_CC_CVV
)

func mkField(label string, size int, fn func(string) (string, bool)) form.Field {
	return form.Label(form.NewTextField(size, []rune(""), fn), label)
}
func mkPasswordFields(size int) (passField, confirmField form.Field) {
	passTextField := form.NewTextField(size, []rune(""), validPassword)
	passField = form.Label(passTextField, "Password")
	confirmTextField := form.NewTextField(size, []rune(""), func(val string) (string, bool) {
		if prob, ok := validPassword(val); !ok {
			return prob, ok
		}
		if val != passTextField.Value() {
			return "Passwords not identical", false
		}
		return "", true

	})
	confirmField = form.Label(confirmTextField, "Enter the password again for confirmation")

	return
}

func MakeSignupForm() (fields map[int]form.Field, f *form.Form, signup *bool) {
	pass, confirm := mkPasswordFields(24)
	fields = map[int]form.Field{
		FIELD_ACCOUNT:            mkField("Account name (this will be used as part of your machines hostnames)", 24, validName),
		FIELD_OWNER_NAME:         mkField("Account owner's username\r\nThis is the name you will use to log in. At a later time you can add a technical contact to the account by emailing support. This tool will eventually have support for that also.", 24, validName),
		FIELD_OWNER_EMAIL:        mkField("Email address", 24, validNonEmpty), // TODO(telyn): make sure it's email-lookin'
		FIELD_OWNER_PASS:         pass,
		FIELD_OWNER_PASS_CONFIRM: confirm,
		FIELD_OWNER_FIRSTNAME:    mkField("First name", 24, validNonEmpty),
		FIELD_OWNER_LASTNAME:     mkField("Last name", 24, validNonEmpty),
		FIELD_OWNER_CC:           mkField("ISO Country code (2-digit country code)\r\nNote that the UK's code is actually GB. Most others are what you'd expect", 3, validISOCountry),
		FIELD_OWNER_POSTCODE:     mkField("Post code", 24, validPostcode),
		FIELD_OWNER_CITY:         mkField("City", 24, validNonEmpty),
		FIELD_OWNER_ADDRESS:      mkField("Street Address", 24, validNonEmpty),
		FIELD_OWNER_PHONE:        mkField("Phone number", 24, validNumber),
		FIELD_OWNER_MOBILE:       mkField("Mobile phone (optional)", 24, validEmptyOr(validNumber)),
		FIELD_OWNER_ORG_NAME:     mkField("Organisation name (optional)", 24, validAlways),
		FIELD_OWNER_ORG_DIVISION: mkField("Organisation division (optional)", 24, validAlways),
		FIELD_OWNER_ORG_VAT:      mkField("VAT Number (optional)", 24, validAlways),
		FIELD_CC_NUMBER:          mkField("Debit/Credit card number", 17, validCC),
		FIELD_CC_NAME:            mkField("Name on card", 17, validNonEmpty),
		FIELD_CC_EXPIRY:          mkField("Expiry (MM/YY)", 6, validExpiry),
		FIELD_CC_CVV:             mkField("CVV2 number (3-4 digit number on back of card)", 5, validCVV),
	}
	fieldsArr := make([]form.Field, len(fields)+2)
	for i, f := range fields {
		fieldsArr[i+1] = f
	}
	fieldsArr[0] = form.NewLabelField("Welcome to Bytemark!\r\n\r\nFilling out this form will create a Bytemark account for you. You can cancel at any time by pressing Esc twice, or using Ctrl+C. Press Tab to cycle through the fields. The fields are underlined in red when invalid, and green when valid.")
	pointer := &f
	s := false
	signup = &s
	fieldsArr[len(fields)+1] = form.NewButtonsField([]form.Button{
		{
			Text: "Sign up",
			Action: func() {
				// this is some fun to prevent a dependency cycle. it is gross.
				if probs, ok := (*pointer).Validate(); !ok {
					log.Debugf(6, "problems with form: %#v", probs)
				} else {
					*signup = true
					(*pointer).Stop()
				}
			},
		},
		{
			Text: "Cancel",
			Action: func() {
				// this is some fun to prevent cyclical dependencies. it is gross
				(*pointer).Stop()
			},
		},
	})
	f = form.NewForm(fieldsArr)
	return
}
