package util

import (
	"bigv.io/client/util/log"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptYesNo provides a y/n prompt. Returns true if the user enters y, false otherwise.
func PromptYesNo(prompt string) bool {
	return Prompt(prompt+" (y/N) ") == "y"
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
	log.Debugf(999, "%s %s", prompt, res)
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
