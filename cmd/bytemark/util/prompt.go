package util

import (
	"fmt"
)

// PromptYesNo provides a y/n prompt. Returns true if the user enters y, false otherwise.
func PromptYesNo(p Prompter, prompt string) bool {
	return p.Prompt(prompt+" (y/N) ") == "y"
}

// PromptfValidate uses prompt as a format string, values as the values for the prompt, then prompts for input.
// The input is then validated using the validation function, and if the input is invalid, it repeats the prompting.
// Returns the valid input once valid input is put in.
func PromptfValidate(p Prompter, valid func(string) (bool, string), prompt string, values ...interface{}) (input string) {
	ok := false

	for !ok {
		input = Promptf(p, prompt, values...)
		ok, input = valid(input)
	}
	return input
}

// PromptValidate prompts for input, validates it. Repeats until the input is actually valid. Returns the valid input.
func PromptValidate(p Prompter, prompt string, valid func(string) bool) (input string) {
	for input = p.Prompt(prompt); !valid(input); input = p.Prompt(prompt) {
		// la la la
		// this is gross
	}
	return input
}

// Promptf formats its arguments with fmt.Sprintf, prompts for input and then returns it.
func Promptf(p Prompter, promptFormat string, values ...interface{}) string {
	return p.Prompt(fmt.Sprintf(promptFormat, values...))
}
