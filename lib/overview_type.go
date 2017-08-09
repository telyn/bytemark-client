package lib

import (
	"io"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// Overview is a combination of a user's default account, their username, and all the accounts they have access to see.
type Overview struct {
	DefaultAccount Account
	Username       string
	Accounts       []Account
}

func (o Overview) DefaultFields(f output.Format) string {
	return "Username, DefaultAccount, Accounts"
}

func (o Overview) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	return FormatOverview(wr, o.Accounts, o.Username)
}
