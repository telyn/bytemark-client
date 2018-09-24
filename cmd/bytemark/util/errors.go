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

// UsageDisplayedError is returned by commands when the user entered wrong info and the help was output
type UsageDisplayedError struct {
	TheProblem string
	Command    string
}

func (e UsageDisplayedError) Error() string {
	return e.TheProblem + "\r\n\r\nFor more information, see `bytemark help " + e.Command + "`"
}

// WontDeleteGroupWithVMsError is returned when 'delete group' was called on a group with stuff in, without --recursive being specified
type WontDeleteGroupWithVMsError struct {
	Group lib.GroupName
}

func (e WontDeleteGroupWithVMsError) Error() string {
	return fmt.Sprintf("Group %s contains servers, will not be deleted without --recursive\r\n", e.Group)
}
