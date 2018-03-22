package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Prompter is an object responsible for prompting the user for input
type Prompter interface {
	Prompt(prompt string) (input string)
}

type realPrompter struct {
	wr io.Writer
	r  io.Reader
}

func (rp realPrompter) Prompt(prompt string) string {
	_, err = fmt.Fprint(rp.wr, prompt)
	if err != nil {
		panic("couldn't prompt. bailing")
	}

	reader := bufio.NewReader(rp.r)
	res, err := reader.ReadString('\n')

	if err != nil {
		if err.Error() == "EOF" {
			return res
		}
		return ""
	}
	return strings.TrimSpace(res)
}

// NewPrompter creates a Prompter which uses stderr and stdin for output and input respectively.
func NewPrompter() Prompter {
	return realPrompter{wr: os.Stderr, r: os.Stdin}
}
