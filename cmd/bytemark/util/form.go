package util

import (
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/telyn/form"
)

const (
	// FormFieldOwnerName is the index of the account owner's username field in the signup form
	FormFieldOwnerName = iota
	// FormFieldOwnerPassword is the index of the account owner's password field in the signup form
	FormFieldOwnerPassword
	// FormFieldOwnerPasswordConfirmation is the index of the account owner's password confirmation field in the signup form
	FormFieldOwnerPasswordConfirmation
	// FormFieldOwnerEmail is the index of the account owner's email address field in the signup form
	FormFieldOwnerEmail
	// FormFieldOwnerFirstName is the index of the account owner's first name field in the signup form
	FormFieldOwnerFirstName
	// FormFieldOwnerLastName is the index of the account owner's last name field in the signup form
	FormFieldOwnerLastName
	// FormFieldOwnerAddress is the index of the account owner's address field in the signup form
	FormFieldOwnerAddress
	// FormFieldOwnerCity is the index of the account owner's city field in the signup form
	FormFieldOwnerCity
	// FormFieldOwnerPostcode is the index of the account owner's postcode field in the signup form
	FormFieldOwnerPostcode
	// FormFieldOwnerCountryCode is the index of the account owner's country code field in the signup form
	FormFieldOwnerCountryCode
	// FormFieldOwnerPhoneNumber is the index of the account owner's phone number field in the signup form
	FormFieldOwnerPhoneNumber
	// FormFieldOwnerMobileNumber is the index of the account owner's mobile number field in the signup form
	FormFieldOwnerMobileNumber
	// FormFieldOwnerOrgName is the index of the account owner's organisation name field in the signup form
	FormFieldOwnerOrgName
	// FormFieldOwnerOrgDivision is the index of the account owner's organisation division field in the signup form
	FormFieldOwnerOrgDivision
	// FormFieldOwnerOrgVATNumber is the index of the account owner's organisation VAT code field in the signup form
	FormFieldOwnerOrgVATNumber
	// FormFieldCreditCardNumber is the index of the credit card number field in the signup form
	FormFieldCreditCardNumber
	// FormFieldCreditCardName is the index of the credit card full name field in the signup form
	FormFieldCreditCardName
	// FormFieldCreditCardExpiry is the index of the credit card expiry field in the signup form
	FormFieldCreditCardExpiry
	// FormFieldCreditCardCVV is the index of the credit card cvv field in the signup form
	FormFieldCreditCardCVV
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
		FormFieldOwnerName:                 mkField("Account name\r\nThis will be the name you use to log in, as well as part of your server's host names.", 24, validName),
		FormFieldOwnerEmail:                mkField("Email address", 24, validNonEmpty), // TODO(telyn): make sure it's email-lookin'
		FormFieldOwnerPassword:             pass,
		FormFieldOwnerPasswordConfirmation: confirm,
		FormFieldOwnerFirstName:            mkField("First name", 24, validNonEmpty),
		FormFieldOwnerLastName:             mkField("Last name", 24, validNonEmpty),
		FormFieldOwnerAddress:              mkField("Street Address", 24, validNonEmpty),
		FormFieldOwnerCity:                 mkField("City", 24, validNonEmpty),
		FormFieldOwnerPostcode:             mkField("Post code", 24, validPostcode),
		FormFieldOwnerCountryCode:          mkField("ISO Country code (2-digit country code)\r\nNote that the UK's code is actually GB. Most others are what you'd expect", 3, validISOCountry),
		FormFieldOwnerPhoneNumber:          mkField("Phone number", 24, validNumber),
		FormFieldOwnerMobileNumber:         mkField("Mobile phone (optional)", 24, validEmptyOr(validNumber)),
		FormFieldOwnerOrgName:              mkField("Organisation name (optional)", 24, validAlways),
		FormFieldOwnerOrgDivision:          mkField("Organisation division (optional)", 24, validAlways),
		FormFieldOwnerOrgVATNumber:         mkField("VAT Number (optional)", 24, validAlways),
	}
	if creditCardForm {
		fields[FormFieldCreditCardNumber] = mkField("Debit/Credit card number", 17, validCC)
		fields[FormFieldCreditCardName] = mkField("Name on card", 17, validNonEmpty)
		fields[FormFieldCreditCardExpiry] = mkField("Expiry (MM/YY)", 6, validExpiry)
		fields[FormFieldCreditCardCVV] = mkField("CVV2 number (3-4 digit number on back of card)", 5, validCVV)
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
