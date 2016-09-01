package util

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/telyn/form"
)

const (
	// FIELD_OWNER_NAME is the index of the account owner's username field in the signup form
	FIELD_OWNER_NAME = iota
	// FIELD_OWNER_PASS is the index of the account owner's password field in the signup form
	FIELD_OWNER_PASS
	// FIELD_OWNER_PASS_CONFIRM is the index of the account owner's password confirmation field in the signup form
	FIELD_OWNER_PASS_CONFIRM
	// FIELD_OWNER_EMAIL is the index of the account owner's email address field in the signup form
	FIELD_OWNER_EMAIL
	// FIELD_OWNER_FIRSTNAME is the index of the account owner's first name field in the signup form
	FIELD_OWNER_FIRSTNAME
	// FIELD_OWNER_LASTNAME is the index of the account owner's last name field in the signup form
	FIELD_OWNER_LASTNAME
	// FIELD_OWNER_CC is the index of the account owner's country code field in the signup form
	FIELD_OWNER_CC
	// FIELD_OWNER_POSTCODE is the index of the account owner's postcode field in the signup form
	FIELD_OWNER_POSTCODE
	// FIELD_OWNER_CITY is the index of the account owner's city field in the signup form
	FIELD_OWNER_CITY
	// FIELD_OWNER_ADDRESS is the index of the account owner's address field in the signup form
	FIELD_OWNER_ADDRESS
	// FIELD_OWNER_PHONE is the index of the account owner's phone number field in the signup form
	FIELD_OWNER_PHONE
	// FIELD_OWNER_MOBILE is the index of the account owner's mobile number field in the signup form
	FIELD_OWNER_MOBILE
	// FIELD_OWNER_ORG_NAME is the index of the account owner's organisation name field in the signup form
	FIELD_OWNER_ORG_NAME
	// FIELD_OWNER_ORG_DIVISION is the index of the account owner's organisation division field in the signup form
	FIELD_OWNER_ORG_DIVISION
	// FIELD_OWNER_ORG_VAT is the index of the account owner's organisation VAT code field in the signup form
	FIELD_OWNER_ORG_VAT
	// FIELD_CC_NUMBER is the index of the credit card number field in the signup form
	FIELD_CC_NUMBER
	// FIELD_CC_NAME is the index of the credit card full name field in the signup form
	FIELD_CC_NAME
	// FIELD_CC_EXPIRY is the index of the credit card expiry field in the signup form
	FIELD_CC_EXPIRY
	// FIELD_CC_CVV is the index of the credit card cvv field in the signup form
	FIELD_CC_CVV
)

func mkField(label string, size int, fn func(string) (string, bool)) form.Field {
	return form.Label(form.NewTextField(size, []rune(""), fn), label)
}
func mkPasswordFields(size int) (passField, confirmField form.Field) {
	passTextField := form.NewMaskedTextField(size, []rune(""), validPassword)
	passField = form.Label(passTextField, "Password")
	confirmTextField := form.NewMaskedTextField(size, []rune(""), func(val string) (string, bool) {
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

// MakeSignupForm constructs the singup form, returning a may of all the fields, the form itself and a pointer to a bool that will be true when the user has requested to sign up, and false otherwise. This is so that once the signup form exits you know whether to continue or not.
func MakeSignupForm(creditCardForm bool) (fields map[int]form.Field, f *form.Form, signup *bool) {
	pass, confirm := mkPasswordFields(24)
	fields = map[int]form.Field{
		FIELD_OWNER_NAME:         mkField("Account name\r\nThis will be the name you use to log in, as well as part of your server's host names.", 24, validName),
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
	}
	if creditCardForm {
		fields[FIELD_CC_NUMBER] = mkField("Debit/Credit card number", 17, validCC)
		fields[FIELD_CC_NAME] = mkField("Name on card", 17, validNonEmpty)
		fields[FIELD_CC_EXPIRY] = mkField("Expiry (MM/YY)", 6, validExpiry)
		fields[FIELD_CC_CVV] = mkField("CVV2 number (3-4 digit number on back of card)", 5, validCVV)
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
