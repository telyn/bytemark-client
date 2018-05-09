package auth

import mock "github.com/maraino/go-mock"

type testPrompter struct {
	mock.Mock
}

func (tp testPrompter) Prompt(prompt string) (response string) {
	r := tp.Called(prompt)
	return r.String(0)
}

func (tp testPrompter) Ask(prompt string) (password string, err error) {
	r := tp.Called(prompt)
	return r.String(0), r.Error(1)
}
