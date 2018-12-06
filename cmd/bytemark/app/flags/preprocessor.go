package flags

import "github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"

// A Preprocesser is a Flag that has a preprocess step that requires a Context
type Preprocesser interface {
	Preprocess(c *app.Context) error
}
