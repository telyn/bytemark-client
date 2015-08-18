package cmd

import (
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
	fmt.Print(prompt)
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
