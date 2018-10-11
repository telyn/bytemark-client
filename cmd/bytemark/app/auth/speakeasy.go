package auth

import (
	"os"

	"github.com/bgentry/speakeasy"
)

type passPrompter interface {
	Ask(prompt string) (password string, err error)
}

type speakeasyWrapper struct{}

func (sw speakeasyWrapper) Ask(prompt string) (string, error) {
	return speakeasy.FAsk(os.Stderr, prompt)
}
