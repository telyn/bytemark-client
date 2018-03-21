package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Prompter interface {
	Prompt(string) string
}

type realPrompter struct {
	wr io.Writer
	r  io.Reader
}

func (rp realPrompter) Prompt(prompt string) string {
	fmt.Fprint(rp.wr, prompt)

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

func NewPrompter() Prompter {
	return realPrompter{wr: os.Stderr, r: os.Stdin}
}
