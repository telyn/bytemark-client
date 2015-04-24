package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PromptYesNo(prompt string) bool {
	return Prompt(prompt+" (y/n) ") == "y"
}

func Prompt(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')

	if err != nil {
		exit(err)
	}
	return strings.TrimSpace(res)
}
