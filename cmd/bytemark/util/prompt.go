package util

import (
	"bufio"
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"os"
	"regexp"
	"strings"
)

// PromptYesNo provides a y/n prompt. Returns true if the user enters y, false otherwise.
func PromptYesNo(prompt string) bool {
	return Prompt(prompt+" (y/N) ") == "y"
}

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

// PromptfValidate uses prompt as a format string, values as the values for the prompt, then prompts for input.
// The input is then validated using the validation function, and if the input is invalid, it repeats the prompting.
// Returns the valid input once valid input is put in.
func PromptfValidate(valid func(string) (bool, string), prompt string, values ...interface{}) (input string) {
	ok := false

	for !ok {
		input = Promptf(prompt, values)
		ok, input = valid(input)
	}
	return input
}

// PromptValidate prompts for input, validates it. Repeats until the input is actually valid. Returns the valid input.
func PromptValidate(prompt string, valid func(string) bool) (input string) {
	for input = Prompt(prompt); !valid(input); input = Prompt(prompt) {
		// la la la
		// this is gross
	}
	return input
}

// Promptf formats its arguments with fmt.Sprintf, prompts for input and then returns it.
func Promptf(promptFormat string, values ...interface{}) string {
	return Prompt(fmt.Sprintf(promptFormat, values...))
}

// Prompt provides a string prompt, returns the entered string with no whitespace
func Prompt(prompt string) string {
	fmt.Fprint(os.Stderr, prompt)
	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')

	if err != nil {
		if err.Error() == "EOF" {
			return res
		}
		return ""
	}
	return strings.TrimSpace(res)
}
