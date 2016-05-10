package util

import (
	"bufio"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// PromptYesNo provides a y/n prompt. Returns true if the user enters y, false otherwise.
func PromptYesNo(prompt string) bool {
	return Prompt(prompt+" (y/N) ") == "y"
}

func ValidCreditCard(input string) (bool, string) {
	r := regexp.MustCompile("/^([0-9]{4} ?){4}$")
	if r.MatchString(input) {
		return true, strings.Replace(input, " ", "", -1)
	}
	log.Log("Credit card numbers should be specified as 16 numbers, with or without spaces")
	// check digit
	return false, input
}

func ValidExpiry(input string) (bool, string) {
	r := regexp.MustCompile("/^[0-9]{4}$/")
	if r.MatchString(input) {
		return true, input
	}
	log.Log("Expiry dates should be in MMYY format without a slash or space.")
	return false, input
}

func ValidEmail(input string) (bool, string) {
	r := regexp.MustCompile("/^.*@([a-z0-9A-Z-]+\\.)+[a-z0-9A-Z-]$/")
	if r.MatchString(input) {
		return true, input
	}
	log.Log("Email address must have a local part and a domain.")
	return false, input
}

func ValidName(input string) (bool, string) {
	if strings.Contains(input, " ") {
		return false, input
	}
	return true, input
}

func PromptfValidate(valid func(string) (bool, string), prompt string, values ...interface{}) (input string) {
	ok := false

	for !ok {
		input = Promptf(prompt, values)
		ok, input = valid(input)
	}
	return input
}

func PromptValidate(prompt string, valid func(string) bool) (input string) {
	for input = Prompt(prompt); !valid(input); input = Prompt(prompt) {
		// la la la
		// this is gross
	}
	return input
}

func Promptf(promptFormat string, values ...interface{}) string {
	return Prompt(fmt.Sprintf(promptFormat, values...))
}

// Prompt provides a string prompt, returns the entered string with no whitespace (hopefully)
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

func ShiftArgument(args *[]string, kindOfThing string) (string, bool) {
	if len(*args) > 0 {
		value := (*args)[0]
		*args = (*args)[1:]
		return value, true
	} else {
		log.Errorf("Not enough arguments. A %s was not specified.\r\n", kindOfThing)
		return "", false
	}
}
