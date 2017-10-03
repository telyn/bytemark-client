package args

import (
	"strings"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
)

// Optional takes a list of flag names. For each flag name it attempts to read the next arg and set the flag with the corresponding name.
// for instance:
// Optional("server", "disc", "size")
// will attempt to read 3 arguments, setting the "server" flag to the first, "disc" to the 2nd, "size" to the third.
func Optional(args ...string) func(c *app.Context) error {
	return func(c *app.Context) error {
		for _, name := range args {
			value, err := c.NextArg()
			if err != nil {
				// if c.NextArg errors that means there aren't more arguments
				// so we just return nil - returning an error would stop the execution of the action.
				return nil
			}
			err = c.Context.Set(name, value)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Join is like Optional, but reads up to n arguments joined with spaces and sets the one named flag.
// if n is not set, reads all the remaining arguments.
func Join(flagName string, n ...int) func(c *app.Context) error {
	return func(c *app.Context) (err error) {
		toRead := len(c.Args())
		if len(n) > 0 {
			toRead = n[0]
		}

		value := make([]string, 0, toRead)
		for i := 0; i < toRead; i++ {
			arg, argErr := c.NextArg()
			if argErr != nil {
				// don't return the error - just means we ran out of arguments to slurp
				break
			}
			value = append(value, arg)
		}
		err = c.Context.Set(flagName, strings.Join(value, " "))
		return

	}
}
