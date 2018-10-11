package mocks

import mock "github.com/maraino/go-mock"

// Prompter implements the cmd/bytemark/util.Prompter interface,
// as well as the cmd/bytemark/app/auth.passwordPrompter interface
// which allows for wrapping speakeasy.Ask calls
type Prompter struct {
	mock.Mock
}

func (tp Prompter) Prompt(prompt string) (response string) {
	r := tp.Called(prompt)
	return r.String(0)
}

func (tp Prompter) Ask(prompt string) (password string, err error) {
	r := tp.Called(prompt)
	return r.String(0), r.Error(1)
}
