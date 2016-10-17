package util

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib"
)

// SubprocessFailedError is returned when a process run by bytemark-client (e.g. open/xdg-open to open a browser) failed
type SubprocessFailedError struct {
	Args     []string
	ExitCode int
	Err      error
}

func (e SubprocessFailedError) Error() string {
	return fmt.Sprintf("Running %s failed - %s", e.Args[0], e.Err.Error())
}

// NotEnoughArgumentsError is returned from command's Actions when there weren't enough arguments to satisfy the requirements. This is no longer used in favour of return c.Help("not enough arguments") now and will be removed when I do a deadcode delinting pass.
type NotEnoughArgumentsError struct {
}

func (e NotEnoughArgumentsError) Error() string {
	return "Not enough arguments passed to the command!"
}

// UsageDisplayedError is returned by commands when the user entered wrong info and the help was output
type UsageDisplayedError struct {
	TheProblem string
	Command    string
}

func (e UsageDisplayedError) Error() string {
	return e.TheProblem + "\r\n\r\nFor more information, see `bytemark help " + e.Command + "`"
}

// WontDeleteNonEmptyGroupError is returned when 'delete group' was called on a group with stuff in, without --recursive being specified
type WontDeleteNonEmptyGroupError struct {
	Group *lib.GroupName
}

func (e WontDeleteNonEmptyGroupError) Error() string {
	return fmt.Sprintf("Group %s contains servers, will not be deleted without --recursive\r\n", e.Group)
}
