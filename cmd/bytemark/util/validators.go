package util

import (
	"regexp"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/util/log"
)

// ValidCreditCard is a credit-card-looking bunch of numbers. Doesn't check the check digit.
func ValidCreditCard(input string) (bool, string) {
	r := regexp.MustCompile(`/^([0-9]{4} ?){4}$`)
	if r.MatchString(input) {
		return true, strings.Replace(input, " ", "", -1)
	}
	log.Log("Credit card numbers should be specified as 16 numbers, with or without spaces")
	// check digit
	return false, input
}

// ValidExpiry checks that the input is a valid credit card expiry, written in MMYY format.
func ValidExpiry(input string) (bool, string) {
	r := regexp.MustCompile("/^[0-9]{4}$/")
	if r.MatchString(input) {
		return true, input
	}
	log.Log("Expiry dates should be in MMYY format without a slash or space.")
	return false, input
}

// ValidEmail checks that the input looks vaguely like an email address. It's very loose, relies on better validation elsewhere.
func ValidEmail(input string) (bool, string) {
	r := regexp.MustCompile(`/^.*@([a-z0-9A-Z-]+\.)+[a-z0-9A-Z-]$/`)
	if r.MatchString(input) {
		return true, input
	}
	log.Log("Email address must have a local part and a domain.")
	return false, input
}

// ValidName checks to see that the input looks like a name. Names can't have spaces in, that's all I know.
func ValidName(input string) (bool, string) {
	if strings.Contains(input, " ") {
		return false, input
	}
	return true, input
}
